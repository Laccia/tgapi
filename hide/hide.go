package hide

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	vault "github.com/hashicorp/vault/api"
)

// type Myclient struct {
// 	cli *vault.Client
// }

// This is the accompanying code for the Developer Quick Start.
// WARNING: Using root tokens is insecure and should never be done in production!
// func VaultClient() {

// }

func WriteSecret(client *vault.Client) { // json уходит вверх ногами
	// Authenticate

	previous, err := os.ReadFile("app/tgProto/session.json")
	if err != nil {
		fmt.Println("cant find file")
	}

	secretData := map[string]interface{}{}
	// json to map
	ss := json.Unmarshal(previous, &secretData)

	if ss != nil {
		fmt.Println("bad try")
	}
	// log.Print(secretData)

	// Write a secret
	_, err = client.KVv2("kv").Put(context.Background(), "my-secret-password", secretData)
	if err != nil {
		fmt.Printf("unable to write secret: %v", err)
	}
}

// fmt.Println("Secret written successfully.")
func ReadSecret(client *vault.Client) { // вроде ок

	// Read a secret from the default mount path for KV v2 in dev mode, "secret"
	secret, err := client.KVv2("kv").Get(context.Background(), "my-secret-password")
	if err != nil {
		log.Fatalf("unable to read secret: %v", err)
	}

	// value, ok := secret.Data["AuthKeyID"].(string)
	// for k, v := range secret.Data {
	// 	fmt.Println("key:", k, "\n", "value:", v)
	// }
	// file, err := json.Marshal(secret)
	// if err != nil {
	// 	fmt.Println("cant write file")
	// }
	file, _ := json.Marshal(secret.Data)
	ssfile, err := os.OpenFile("app/tgProto/session.json", os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println("cant open file")
	}
	err = os.WriteFile("app/tgProto/session.json", file, 0644)
	if err != nil {
		fmt.Println(err)
	}
	ssfile.Close()

	// if value != "Hashi123" {
	// 	log.Fatalf("unexpected password value %q retrieved from vault", value)
	// }

	fmt.Println("Vault: Access granted!")
}
