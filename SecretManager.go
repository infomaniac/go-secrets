package secman

import (
	"context"
	"fmt"

	gsm "cloud.google.com/go/secretmanager/apiv1"
	gsmt "cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
)

type SecretManager struct {
	ctx       context.Context
	client    *gsm.Client
	projectID string
	locations []string
}

func New(ctx context.Context, projectID string, secretLocations ...string) (*SecretManager, error) {
	client, err := gsm.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &SecretManager{
		ctx:       ctx,
		client:    client,
		projectID: projectID,
		locations: secretLocations,
	}, nil
}
func (s *SecretManager) Close() error {
	return s.client.Close()
}

func (s *SecretManager) Get(secretID string) ([]byte, error) {
	accessRequest := &gsmt.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/latest", s.projectID, secretID),
	}

	result, err := s.client.AccessSecretVersion(s.ctx, accessRequest)
	if err != nil {
		return nil, err
	}

	return result.Payload.Data, nil
}

func (s *SecretManager) Set(secretName string, secretValue []byte) error {
	secret, err := s.getSecret(secretName)
	if secret == nil {
		secret, err = s.createSecret(secretName)
	}
	if err != nil {
		return err
	}

	return s.addSecretVersion(secret, secretValue)
}

func (s *SecretManager) getSecret(secretID string) (*gsmt.Secret, error) {
	req := &gsmt.GetSecretRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s", s.projectID, secretID),
	}
	return s.client.GetSecret(s.ctx, req)
}

func (s *SecretManager) createSecret(secretID string) (*gsmt.Secret, error) {
	var rep *gsmt.Replication

	if len(s.locations) == 0 {
		rep = &gsmt.Replication{
			Replication: &gsmt.Replication_Automatic_{Automatic: &gsmt.Replication_Automatic{}},
		}
	} else {
		replicas := make([]*gsmt.Replication_UserManaged_Replica, 0, len(s.locations))
		for _, loc := range s.locations {
			replicas = append(replicas, &gsmt.Replication_UserManaged_Replica{Location: loc})
		}

		rep = &gsmt.Replication{
			Replication: &gsmt.Replication_UserManaged_{
				UserManaged: &gsmt.Replication_UserManaged{
					Replicas: replicas,
				},
			},
		}
	}

	createSecretReq := &gsmt.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s", s.projectID),
		SecretId: secretID,
		Secret:   &gsmt.Secret{Replication: rep},
	}

	return s.client.CreateSecret(s.ctx, createSecretReq)
}

func (s *SecretManager) addSecretVersion(sec *gsmt.Secret, secretValue []byte) error {
	addSecretVersionReq := &gsmt.AddSecretVersionRequest{
		Parent: sec.Name,
		Payload: &gsmt.SecretPayload{
			Data: secretValue,
		},
	}

	_, err := s.client.AddSecretVersion(s.ctx, addSecretVersionReq)
	return err
}

func (s *SecretManager) Delete(secretID string) error {
	req := &gsmt.DeleteSecretRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s", s.projectID, secretID),
	}

	return s.client.DeleteSecret(s.ctx, req)
}
