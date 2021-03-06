package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var artifactoryNamespace string
var artifactoryCredentialsSecret string
var artifactoryTokenScope string
var artifactoryTokenUserPrefix string
var buildNamespaces namespaces = []string{"build"}
var createDockerRegistrySecret bool
var dockerServer string
var secretName string
var secretKey string
var secretUserKey string

type namespaces []string

type createTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type dockerConfigJSON struct {
	Auths map[string]auth
}

type auth struct {
	Username string
	Password string
}

var clientset *kubernetes.Clientset
var artifactoryUsername string
var artifactoryPassword string

const tokenEndpoint = "/artifactory/api/security/token"
const validityEndpoint = "/artifactory/api/repositories"

func (i *namespaces) Set(value string) error {
	*i = strings.Split(value, ",")
	return nil
}

func (i *namespaces) String() string {
	return fmt.Sprint(*i)
}

func initFlags() {
	flag.StringVar(&artifactoryNamespace, "artifactoryNamespace", "default", "namespace to look for artifactory instance")
	flag.StringVar(&artifactoryCredentialsSecret, "artifactoryCredentialsSecret", "artifactory-credentials", "artifactory admin credentials secret name")
	flag.StringVar(&artifactoryTokenUserPrefix, "artifactoryTokenUserPrefix", "gitlab-", "user prefix for artifactory token")
	flag.StringVar(&artifactoryTokenScope, "artifactoryTokenScope", "", "comma separated groups for artifactory token")
	flag.StringVar(&secretName, "secretName", "artifactory-access-token", "name of the secret containing the token or the docker credentials")
	flag.StringVar(&secretKey, "secretKey", "artifactory-access-token", "key in the secret containing the token (ignored if createDockerRegistrySecret=true)")
	flag.StringVar(&secretUserKey, "secretUserKey", "artifactory-access-token-user", "key in the secret containing the name of the ephemeral user of the token (ignored if createDockerRegistrySecret=true)")
	flag.BoolVar(&createDockerRegistrySecret, "createDockerRegistrySecret", false, "if you want to create a registry credential secret, instead of a normal access-token")
	flag.StringVar(&dockerServer, "dockerServer", "", "url of the docker server")
	flag.Var(&buildNamespaces, "buildNamespaces", "comma separated ci build namespaces to monitor")
	flag.Parse()
}

func main() {
	initFlags()
	log.Println("ci build namespaces to monitor=", buildNamespaces)
	log.Println("namespace to look for artifactory instance=", artifactoryNamespace)
	log.Println("name of the secret containing artifactory access token=", secretName)
	log.Println("name of the key in the secret=", secretKey)

	config, err := rest.InClusterConfig()

	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	watchArtifactory()
}

func upsertAccessTokenSecret(artifactoryURL string) {
	for _, n := range buildNamespaces {
		current, err := clientset.CoreV1().Secrets(n).Get(secretName, metav1.GetOptions{})

		if err != nil && err.(*errors.StatusError).ErrStatus.Code == 404 {
			log.Println("create a new artifactory token in build domain=", n)
			newToken := getNewToken(artifactoryURL, n)
			if newToken != nil {
				_, err := clientset.CoreV1().Secrets(n).Create(newToken)
				if err != nil {
					log.Println("error creating the token", err)
				}
			}
		} else if err != nil {
			log.Println("unespected error occurred", err)
			continue
		} else {
			token, found := getTokenFromSecret(current)
			updateTokenIfNotValid(artifactoryURL, n, token, found)
		}
	}
}

func getTokenFromSecret(currentSecret *v1.Secret) (string, bool) {
	if createDockerRegistrySecret {
		if dockerConfigEntry, found := currentSecret.Data[".dockerconfigjson"]; found {
			return getTokenFromDockerConfigSecret(dockerConfigEntry)
		}
		return "", false
	}	

	if token, found := currentSecret.Data[secretKey]; found {
		return string(token), true
	}
	return "", false
}

func getTokenFromDockerConfigSecret(content []byte) (string, bool) {
	var config dockerConfigJSON
	json.Unmarshal(content, &config)
	if entry, found := config.Auths[dockerServer]; found {
		return entry.Password, true
	}
	return "", false
}

func makeTokenRequest(artifactoryURL string, path string, method string, token string, body io.Reader) (*http.Response, error) {
	u, _ := url.ParseRequestURI(artifactoryURL)
	u.Path = path
	client := &http.Client{}
	r, _ := http.NewRequest(method, u.String(), body)
	addAuthorizationHeader(r, token)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return client.Do(r)
}

func getNewToken(artifactoryURL string, namespace string) *v1.Secret {
	data := url.Values{}
	username := artifactoryTokenUserPrefix + namespace
	data.Set("username", username)
	data.Set("scope", "member-of-groups:\""+artifactoryTokenScope+"\"")
	data.Set("expires_in", "0")
	resp, err := makeTokenRequest(artifactoryURL, tokenEndpoint, "POST", "", strings.NewReader(data.Encode()))
	if err != nil {
		log.Println("cannot contact artifactory service. cannot get a new token. token will not be updated nor created. error=", err)
		return nil
	}
	defer resp.Body.Close()
	tokenResp := createTokenResponse{}
	json.NewDecoder(resp.Body).Decode(&tokenResp)

	if createDockerRegistrySecret {
		return &v1.Secret{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Secret",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: namespace,
			},
			StringData: map[string]string{
				".dockerconfigjson": "{\"auths\":{\"" + dockerServer + "\":{\"username\":\"" + username + "\",\"password\":\"" + tokenResp.AccessToken + "\"}}}",
			},
			Type: "kubernetes.io/dockerconfigjson",
		}
	}
	return &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		StringData: map[string]string{
			secretKey:     tokenResp.AccessToken,
			secretUserKey: username,
		},
	}
}

func addAuthorizationHeader(r *http.Request, token string) {
	if token == "" {
		//use basic auth with username and password
		r.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(artifactoryUsername+":"+artifactoryPassword)))
	} else {
		//use bearer
		r.Header.Add("Authorization", "Bearer "+token)
	}
}

func updateTokenIfNotValid(artifactoryURL string, namespace string, token string, found bool) {
	if found && isTokenValid(artifactoryURL, namespace, token) {
		//logging is done in the above method. just returning
		return
	}
	log.Println("updated invalid artifactory token in build domain=", namespace)
	newToken := getNewToken(artifactoryURL, namespace)
	if newToken != nil {
		//updating invalid secret with new token
		_, err := clientset.CoreV1().Secrets(namespace).Update(newToken)
		if err != nil {
			log.Println("error updating the token", err)
		}
	}
}

func isTokenValid(artifactoryURL string, namespace string, token string) bool {
	resp, err := makeTokenRequest(artifactoryURL, validityEndpoint, "GET", token, nil)
	if err != nil {
		log.Println("cannot contact artifactory service. cannot check token validity. token will not be updated nor created. error=", err)
		//token is marked as valid (even if we don't know), so that the updateTokenIfNotValid caller can return
		return true
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		//token is invalid
		return false
	}
	log.Println("nothing to update. valid artifactory token in build domain=", namespace)
	return true
}

func parseArtifactoryAdminCredentials() {
	//retrieve the secret where artifactory admin credentials are stored
	credentials, err := clientset.CoreV1().Secrets(artifactoryNamespace).Get(artifactoryCredentialsSecret, metav1.GetOptions{})

	if err != nil {
		panic(err.Error())
	}
	artifactoryUsername = string(credentials.Data["username"])
	artifactoryPassword = string(credentials.Data["password"])
}

func watchArtifactory() {
	parseArtifactoryAdminCredentials()
	options := metav1.ListOptions{
		LabelSelector: "app=artifactory",
	}

	for {
		watcher, err := clientset.CoreV1().Endpoints(artifactoryNamespace).Watch(options)
		if err != nil {
			panic(err.Error())
		}
		ch := watcher.ResultChan()
		for event := range ch {
			endpoints := event.Object.(*v1.Endpoints)
			switch event.Type {
			case watch.Added:
				handleModified(endpoints)
			case watch.Deleted:
				//handleDeleted(endpoints)
			case watch.Modified:
				handleModified(endpoints)
			}
		}
	}
}

// func handleDeleted(endpoints *v1.Endpoints) {
// 	log.Println("deleted artifactory")
// 	for _, n := range buildNamespaces {
// 		log.Println("deleting gitlab token in build domain ", n)
// 		clientset.CoreV1().Secrets(n).Delete(secretName, &metav1.DeleteOptions{})
// 	}
// }

func handleModified(endpoints *v1.Endpoints) {
	if len(endpoints.Subsets) < 1 {
		log.Println("found non ready artifactory")
		return
	}
	//monitoring added or modified subsets
	for _, subsets := range endpoints.Subsets {
		if subsets.Addresses != nil && len(subsets.Addresses) > 0 {
			//found ready artifactory address
			hostPort := fmt.Sprintf("%s:%d", subsets.Addresses[0].IP, endpoints.Subsets[0].Ports[0].Port)
			conn, err := net.DialTimeout("tcp", hostPort, 10*time.Second)
			if err != nil {
				log.Println("artifactory service is still unreachable after timeout. processing skipped", err)
				return
			}
			defer conn.Close()
			url := "http://" + hostPort
			log.Println("found ready artifactory ", url)
			upsertAccessTokenSecret(url)
			return
		}
	}
	log.Println("found non ready artifactory")
}
