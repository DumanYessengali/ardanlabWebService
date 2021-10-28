package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

/*
	openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -pubout -in private.pem -out public.pem
*/
func main() {
	//keygen()
	tokengen()
}

func tokengen() {
	privatePEM, err := ioutil.ReadFile("c:/Users/balmu/go/src/github.com/DumanYessengali/ardanlabWebService/private.pem")
	if err != nil {
		log.Fatalln(err)
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privatePEM)
	if err != nil {
		log.Fatalln(err)
	}

	claims := struct {
		jwt.StandardClaims
		Authorized []string
	}{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "service project",
			Subject:   "123456789",
			ExpiresAt: time.Now().Add(8760 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		Authorized: []string{"ADMIN"},
	}

	method := jwt.GetSigningMethod("RS256")
	tkn := jwt.NewWithClaims(method, claims)
	tkn.Header["kid"] = "akjfasndio24noifwe-fwejfn-wefnwe-nlnknkl"
	str, err := tkn.SignedString(privateKey)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(str)

}

func keygen() {
	//private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalln(err)
	}

	privateFile, err := os.Create("private.pem")
	if err != nil {
		log.Fatalln(err)
	}
	defer privateFile.Close()

	privateBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	if err := pem.Encode(privateFile, &privateBlock); err != nil {
		log.Fatalln(err)
	}

	//public key

	ans1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		log.Fatalln(err)
	}

	publicFile, err := os.Create("public.pem")
	if err != nil {
		log.Fatalln(err)
	}
	defer publicFile.Close()

	publicBlock := pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: ans1Bytes,
	}

	if err := pem.Encode(publicFile, &publicBlock); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("DONE")
}
