package main

import "testing"

func Test_getTokenFromDockerConfigSecret(t *testing.T) {
	dockerServer = "https://artifactory.eidosmedia.sh"
	type args struct {
		content []byte
	}
	tests := []struct {
		name string
		args args
		token string
		found bool
	}{
		{
			"Test",
			args{
				[]byte(`{"auths":{"https://artifactory.eidosmedia.sh":{"username":"gitlab-kubernetes-docker-credentials-default","password":"eyJ2ZXIiOiIyIiwidHlwIjoiSldUIiwiYWxnIjoiUlMyNTYiLCJraWQiOiJxZWF2VDJwWVNZZm1FU2c1bGtkNFRJM243dUpCRkw1ZkswaURNVURpRzNvIn0.eyJzdWIiOiJqZnJ0QDAxZHQ0eHE1ZXRzeDI1MXhqeDBua3IxZHp3XC91c2Vyc1wvZ2l0bGFiLWt1YmVybmV0ZXMtZG9ja2VyLWNyZWRlbnRpYWxzLWRlZmF1bHQiLCJzY3AiOiJtZW1iZXItb2YtZ3JvdXBzOlwiZG9ja2VyLWRlcGxveWVyLHJlYWRlcnNcIiBhcGk6KiIsImF1ZCI6ImpmcnRAMDFkdDR4cTVldHN4MjUxeGp4MG5rcjFkenciLCJpc3MiOiJqZnJ0QDAxZHQ0eHE1ZXRzeDI1MXhqeDBua3IxZHp3XC91c2Vyc1wvYWRtaW4iLCJpYXQiOjE1ODE1MDU3MTksImp0aSI6ImUwYTM5ZmM5LWVkMTMtNDg3Yi04OWEyLTEyOGNhYmMwOWIyZSJ9.C4ahzEKPCOUg2YKFx-4bZ3EQXvq7rzQL3axB_1Jmpf3eMMJC-rhGT8rmIyZwSjrVTZWJ9LtOPMU8PD4xIoFRsPDABVghFKlKfiN-1hjajsX0ohFAhC-p_DzzTc0W2HORE2p8Mnsl7Yv6jhkgoWL7O8JYEfstBVkdFjkkV084nyVdWG_i4soCHg-94cA4jtb1TFOVli50tV7FHi8vVq7Frm7QPdSqtQRS7fPzP3q88B0nHr1eIquUBQRAO8DEuPjey1wje7LYuI7YgLYTF6G0gNJhrb56s5zFTDC9G_fNNwVgHbYVSkzvoaXCSL-KREd5IJ4RN_G_9xe36sOK3yycDQ"}}}`),
			},
			"eyJ2ZXIiOiIyIiwidHlwIjoiSldUIiwiYWxnIjoiUlMyNTYiLCJraWQiOiJxZWF2VDJwWVNZZm1FU2c1bGtkNFRJM243dUpCRkw1ZkswaURNVURpRzNvIn0.eyJzdWIiOiJqZnJ0QDAxZHQ0eHE1ZXRzeDI1MXhqeDBua3IxZHp3XC91c2Vyc1wvZ2l0bGFiLWt1YmVybmV0ZXMtZG9ja2VyLWNyZWRlbnRpYWxzLWRlZmF1bHQiLCJzY3AiOiJtZW1iZXItb2YtZ3JvdXBzOlwiZG9ja2VyLWRlcGxveWVyLHJlYWRlcnNcIiBhcGk6KiIsImF1ZCI6ImpmcnRAMDFkdDR4cTVldHN4MjUxeGp4MG5rcjFkenciLCJpc3MiOiJqZnJ0QDAxZHQ0eHE1ZXRzeDI1MXhqeDBua3IxZHp3XC91c2Vyc1wvYWRtaW4iLCJpYXQiOjE1ODE1MDU3MTksImp0aSI6ImUwYTM5ZmM5LWVkMTMtNDg3Yi04OWEyLTEyOGNhYmMwOWIyZSJ9.C4ahzEKPCOUg2YKFx-4bZ3EQXvq7rzQL3axB_1Jmpf3eMMJC-rhGT8rmIyZwSjrVTZWJ9LtOPMU8PD4xIoFRsPDABVghFKlKfiN-1hjajsX0ohFAhC-p_DzzTc0W2HORE2p8Mnsl7Yv6jhkgoWL7O8JYEfstBVkdFjkkV084nyVdWG_i4soCHg-94cA4jtb1TFOVli50tV7FHi8vVq7Frm7QPdSqtQRS7fPzP3q88B0nHr1eIquUBQRAO8DEuPjey1wje7LYuI7YgLYTF6G0gNJhrb56s5zFTDC9G_fNNwVgHbYVSkzvoaXCSL-KREd5IJ4RN_G_9xe36sOK3yycDQ",
			true,
		},
		{
			"Test invalid url",
			args{
				[]byte(`{"auths":{"https://artifactory.eidosmedia.com":{"username":"gitlab-kubernetes-docker-credentials-default","password":"eyJ2ZXIiOiIyIiwidHlwIjoiSldUIiwiYWxnIjoiUlMyNTYiLCJraWQiOiJxZWF2VDJwWVNZZm1FU2c1bGtkNFRJM243dUpCRkw1ZkswaURNVURpRzNvIn0.eyJzdWIiOiJqZnJ0QDAxZHQ0eHE1ZXRzeDI1MXhqeDBua3IxZHp3XC91c2Vyc1wvZ2l0bGFiLWt1YmVybmV0ZXMtZG9ja2VyLWNyZWRlbnRpYWxzLWRlZmF1bHQiLCJzY3AiOiJtZW1iZXItb2YtZ3JvdXBzOlwiZG9ja2VyLWRlcGxveWVyLHJlYWRlcnNcIiBhcGk6KiIsImF1ZCI6ImpmcnRAMDFkdDR4cTVldHN4MjUxeGp4MG5rcjFkenciLCJpc3MiOiJqZnJ0QDAxZHQ0eHE1ZXRzeDI1MXhqeDBua3IxZHp3XC91c2Vyc1wvYWRtaW4iLCJpYXQiOjE1ODE1MDU3MTksImp0aSI6ImUwYTM5ZmM5LWVkMTMtNDg3Yi04OWEyLTEyOGNhYmMwOWIyZSJ9.C4ahzEKPCOUg2YKFx-4bZ3EQXvq7rzQL3axB_1Jmpf3eMMJC-rhGT8rmIyZwSjrVTZWJ9LtOPMU8PD4xIoFRsPDABVghFKlKfiN-1hjajsX0ohFAhC-p_DzzTc0W2HORE2p8Mnsl7Yv6jhkgoWL7O8JYEfstBVkdFjkkV084nyVdWG_i4soCHg-94cA4jtb1TFOVli50tV7FHi8vVq7Frm7QPdSqtQRS7fPzP3q88B0nHr1eIquUBQRAO8DEuPjey1wje7LYuI7YgLYTF6G0gNJhrb56s5zFTDC9G_fNNwVgHbYVSkzvoaXCSL-KREd5IJ4RN_G_9xe36sOK3yycDQ"}}}`),
			},
			"",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if token, found := getTokenFromDockerConfigSecret(tt.args.content); token != tt.token || found != tt.found {
				t.Errorf("getTokenFromDockerConfigSecret() = %v, want %v", token, tt.token)
			}
		})
	}
}
