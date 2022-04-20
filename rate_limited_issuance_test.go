package pat

import (
	"bytes"
	"crypto"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/sha512"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	hpke "github.com/cisco/go-hpke"
	"golang.org/x/crypto/cryptobyte"

	"github.com/cloudflare/pat-go/ecdsa"
	"github.com/cloudflare/pat-go/ed25519"
)

// 4096-bit RSA private key
const testTokenPrivateKey = `
-----BEGIN RSA PRIVATE KEY-----
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
-----END RSA PRIVATE KEY-----`

const (
	outputAnonOriginIDTestVectorEnvironmentKey = "RATE_LIMITED_ANON_ORIGIN_ID_TEST_VECTORS_OUT"
	inputAnonOriginIDTestVectorEnvironmentKey  = "RATE_LIMITED_ANON_ORIGIN_ID_TEST_VECTORS_IN"

	outputOriginEncryptionTestVectorEnvironmentKey = "RATE_LIMITED_ORIGIN_ENCRYPTION_TEST_VECTORS_OUT"
	inputOriginEncryptionTestVectorEnvironmentKey  = "RATE_LIMITED_ORIGIN_ENCRYPTION_TEST_VECTORS_IN"
)

func loadPrivateKey(t *testing.T) *rsa.PrivateKey {
	block, _ := pem.Decode([]byte(testTokenPrivateKey))
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		t.Fatal("PEM private key decoding failed")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		t.Fatal(err)
	}

	return privateKey
}

func TestSignatureDifferences(t *testing.T) {
	_, secretKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Error(err)
	}

	message := make([]byte, 32)
	rand.Reader.Read(message)

	blind := make([]byte, 32)
	rand.Reader.Read(blind)
	signature1 := ed25519.BlindKeySign(secretKey, message, blind)

	rand.Reader.Read(blind)
	signature2 := ed25519.BlindKeySign(secretKey, message, blind)

	if bytes.Equal(signature1[:32], signature2[:32]) {
		t.Fatal("Signature prefix matched when it should vary")
	}
	if bytes.Equal(signature1[32:], signature2[32:]) {
		t.Fatal("Signature prefix matched when it should vary")
	}
}

func TestRateLimitedIssuanceRoundTrip(t *testing.T) {
	issuer := NewRateLimitedIssuer(loadPrivateKey(t))
	testOrigin := "origin.example"
	issuer.AddOrigin(testOrigin)

	curve := elliptic.P384()
	secretKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	blindKey, err := ecdsa.GenerateKey(curve, rand.Reader)
	client := CreateRateLimitedClientFromSecret(secretKey.D.Bytes())

	challenge := make([]byte, 32)
	rand.Reader.Read(challenge)

	nonce := make([]byte, 32)
	rand.Reader.Read(nonce)

	tokenKeyID := issuer.TokenKeyID()
	tokenPublicKey := issuer.TokenKey()
	originIndexKey := issuer.OriginIndexKey(testOrigin)

	requestState, err := client.CreateTokenRequest(challenge, nonce, blindKey.D.Bytes(), tokenKeyID, tokenPublicKey, testOrigin, issuer.NameKey())
	if err != nil {
		t.Error(err)
	}

	blindedSignature, blindedPublicKey, err := issuer.Evaluate(requestState.Request())
	if err != nil {
		t.Error(err)
	}

	publicKeyEnc := elliptic.MarshalCompressed(curve, client.secretKey.PublicKey.X, client.secretKey.PublicKey.Y)

	expectedIndexKey, err := ecdsa.BlindPublicKey(curve, &client.secretKey.PublicKey, originIndexKey)
	if err != nil {
		t.Error(err)
	}
	expectedIndexKeyEnc := elliptic.MarshalCompressed(curve, expectedIndexKey.X, expectedIndexKey.Y)

	expectedIndex, err := computeIndex(publicKeyEnc, expectedIndexKeyEnc)
	if err != nil {
		t.Error(err)
	}

	index, err := FinalizeIndex(publicKeyEnc, blindKey.D.Bytes(), blindedPublicKey)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(index, expectedIndex) {
		t.Fatal("index computation incorrect")
	}

	token, err := requestState.FinalizeToken(blindedSignature)
	if err != nil {
		t.Error(err)
	}

	b := cryptobyte.NewBuilder(nil)
	b.AddUint16(RateLimitedTokenType)
	b.AddBytes(nonce)
	context := sha256.Sum256(challenge)
	b.AddBytes(context[:])
	b.AddBytes(tokenKeyID)
	tokenInput := b.BytesOrPanic()

	hash := sha512.New384()
	hash.Write(tokenInput)
	digest := hash.Sum(nil)
	err = rsa.VerifyPSS(tokenPublicKey, crypto.SHA384, digest, token.Authenticator, &rsa.PSSOptions{
		Hash:       crypto.SHA384,
		SaltLength: crypto.SHA384.Size(),
	})
	if err != nil {
		t.Error(err)
	}
}

///////
// Infallible Serialize / Deserialize
func fatalOnError(t *testing.T, err error, msg string) {
	realMsg := fmt.Sprintf("%s: %v", msg, err)
	if err != nil {
		if t != nil {
			t.Fatalf(realMsg)
		} else {
			panic(realMsg)
		}
	}
}

func mustUnhex(t *testing.T, h string) []byte {
	out, err := hex.DecodeString(h)
	fatalOnError(t, err, "Unhex failed")
	return out
}

func mustHex(d []byte) string {
	return hex.EncodeToString(d)
}

///////
// Index computation test vector structure
type rawOriginEncryptionTestVector struct {
	KEMID               hpke.KEMID  `json:"kem_id"`
	KDFID               hpke.KDFID  `json:"kdf_id"`
	AEADID              hpke.AEADID `json:"aead_id"`
	OriginNameKeySeed   string      `json:"origin_name_key_seed"`
	OriginNameKey       string      `json:"origin_name_key"`
	TokenType           uint16      `json:"token_type"`
	OriginNameKeyID     string      `json:"origin_name_key_id"`
	IndexRequest        string      `json:"request_key"`
	TokenKeyID          uint8       `json:"token_key_id"`
	BlindMessage        string      `json:"blinded_msg"`
	OriginName          string      `json:"origin_name"`
	EncryptedOriginName string      `json:"encrypted_origin_name"`
}

type originEncryptionTestVector struct {
	t                   *testing.T
	kemID               hpke.KEMID
	kdfID               hpke.KDFID
	aeadID              hpke.AEADID
	nameKeySeed         []byte
	nameKey             PrivateNameKey
	tokenType           uint16
	indexRequest        []byte
	tokenKeyID          uint8
	blindMessage        []byte
	issuerKeyID         []byte
	originName          string
	encryptedOriginName []byte
}

type originEncryptionTestVectorArray struct {
	t       *testing.T
	vectors []originEncryptionTestVector
}

func (tva originEncryptionTestVectorArray) MarshalJSON() ([]byte, error) {
	return json.Marshal(tva.vectors)
}

func (tva *originEncryptionTestVectorArray) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &tva.vectors)
	if err != nil {
		return err
	}

	for i := range tva.vectors {
		tva.vectors[i].t = tva.t
	}
	return nil
}

func (etv originEncryptionTestVector) MarshalJSON() ([]byte, error) {
	return json.Marshal(rawOriginEncryptionTestVector{
		KEMID:               etv.kemID,
		KDFID:               etv.kdfID,
		AEADID:              etv.aeadID,
		OriginNameKeySeed:   mustHex(etv.nameKeySeed),
		OriginNameKey:       mustHex(etv.nameKey.Public().Marshal()),
		TokenType:           etv.tokenType,
		IndexRequest:        mustHex(etv.indexRequest),
		TokenKeyID:          etv.tokenKeyID,
		BlindMessage:        mustHex(etv.blindMessage),
		OriginNameKeyID:     mustHex(etv.issuerKeyID),
		OriginName:          mustHex([]byte(etv.originName)),
		EncryptedOriginName: mustHex(etv.encryptedOriginName),
	})
}

func (etv *originEncryptionTestVector) UnmarshalJSON(data []byte) error {
	raw := rawOriginEncryptionTestVector{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	etv.kemID = raw.KEMID
	etv.kdfID = raw.KDFID
	etv.aeadID = raw.AEADID
	etv.nameKeySeed = mustUnhex(nil, raw.OriginNameKeySeed)
	etv.tokenType = raw.TokenType
	if etv.tokenType != RateLimitedTokenType {
		return fmt.Errorf("Unsupported token type")
	}

	if raw.KEMID != hpke.DHKEM_X25519 ||
		raw.KDFID != hpke.KDF_HKDF_SHA256 ||
		raw.AEADID != hpke.AEAD_AESGCM128 {
		// Unsupported ciphersuite -- pass
		return fmt.Errorf("Unsupported ciphersuite")
	}

	nameKey, err := CreatePrivateNameKeyFromSeed(etv.nameKeySeed)
	if err != nil {
		return err
	}
	etv.nameKey = nameKey
	etv.indexRequest = mustUnhex(nil, raw.IndexRequest)
	etv.tokenKeyID = raw.TokenKeyID
	etv.blindMessage = mustUnhex(nil, raw.BlindMessage)
	etv.issuerKeyID = mustUnhex(nil, raw.OriginNameKeyID)
	etv.originName = string(mustUnhex(nil, raw.OriginName))
	etv.encryptedOriginName = mustUnhex(nil, raw.EncryptedOriginName)

	return nil
}

func generateOriginEncryptionTestVector(t *testing.T, kemID hpke.KEMID, kdfID hpke.KDFID, aeadID hpke.AEADID) originEncryptionTestVector {
	ikm := make([]byte, 32)
	rand.Reader.Read(ikm)
	nameKey, err := CreatePrivateNameKeyFromSeed(ikm)
	if err != nil {
		t.Fatal(err)
	}

	// Generate random token and index requests
	indexRequest := make([]byte, 49)
	rand.Reader.Read(indexRequest)
	tokenKeyIDBuf := []byte{0x00}
	rand.Reader.Read(tokenKeyIDBuf)
	blindMessage := make([]byte, 512)
	rand.Reader.Read(blindMessage)

	originName := "test.example"
	_, encryptedName, err := encryptOriginName(nameKey.Public(), tokenKeyIDBuf[0], blindMessage, indexRequest, originName)
	if err != nil {
		t.Fatal(err)
	}

	issuerKeyEnc := nameKey.Public().Marshal()
	issuerKeyID := sha256.Sum256(issuerKeyEnc)

	return originEncryptionTestVector{
		kemID:               kemID,
		kdfID:               kdfID,
		aeadID:              aeadID,
		nameKeySeed:         ikm,
		nameKey:             nameKey,
		issuerKeyID:         issuerKeyID[:],
		tokenType:           RateLimitedTokenType,
		indexRequest:        indexRequest,
		tokenKeyID:          tokenKeyIDBuf[0],
		blindMessage:        blindMessage,
		originName:          originName,
		encryptedOriginName: encryptedName,
	}
}

func verifyOriginEncryptionTestVector(t *testing.T, vector originEncryptionTestVector) {
	suite, err := hpke.AssembleCipherSuite(vector.kemID, vector.kdfID, vector.aeadID)
	if err != nil {
		t.Fatal(err)
	}

	if suite.KEM.ID() != hpke.DHKEM_X25519 ||
		suite.KDF.ID() != hpke.KDF_HKDF_SHA256 ||
		suite.AEAD.ID() != hpke.AEAD_AESGCM128 {
		// Unsupported ciphersuite -- pass
		return
	}

	privateNameKey, err := CreatePrivateNameKeyFromSeed(vector.nameKeySeed)
	if err != nil {
		t.Fatal(err)
	}

	originName, err := decryptOriginName(privateNameKey, vector.tokenKeyID, vector.blindMessage, vector.indexRequest, vector.encryptedOriginName)
	if err != nil {
		t.Fatal(err)
	}

	if originName != vector.originName {
		t.Fatal("origin decryption mismatch")
	}
}

func verifyOriginEncryptionTestVectors(t *testing.T, encoded []byte) {
	vectors := originEncryptionTestVectorArray{t: t}
	err := json.Unmarshal(encoded, &vectors)
	if err != nil {
		t.Fatalf("Error decoding test vector string: %v", err)
	}

	for _, vector := range vectors.vectors {
		verifyOriginEncryptionTestVector(t, vector)
	}
}

func TestVectorGenerateOriginEncryption(t *testing.T) {
	vectors := make([]originEncryptionTestVector, 0)
	vectors = append(vectors, generateOriginEncryptionTestVector(t, hpke.DHKEM_X25519, hpke.KDF_HKDF_SHA256, hpke.AEAD_AESGCM128))

	// Encode the test vectors
	encoded, err := json.Marshal(vectors)
	if err != nil {
		t.Fatalf("Error producing test vectors: %v", err)
	}

	// Verify that we process them correctly
	verifyOriginEncryptionTestVectors(t, encoded)

	var outputFile string
	if outputFile = os.Getenv(outputOriginEncryptionTestVectorEnvironmentKey); len(outputFile) > 0 {
		err := ioutil.WriteFile(outputFile, encoded, 0644)
		if err != nil {
			t.Fatalf("Error writing test vectors: %v", err)
		}
	}
}

func TestVectorVerifyOriginEncryption(t *testing.T) {
	var inputFile string
	if inputFile = os.Getenv(inputOriginEncryptionTestVectorEnvironmentKey); len(inputFile) == 0 {
		t.Skip("Test vectors were not provided")
	}

	encoded, err := ioutil.ReadFile(inputFile)
	if err != nil {
		t.Fatalf("Failed reading test vectors: %v", err)
	}

	verifyOriginEncryptionTestVectors(t, encoded)
}

///////
// Index computation test vector structure
type rawAnonOriginIDTestVector struct {
	ClientSecret  string `json:"sk_sign"`
	ClientPublic  string `json:"pk_sign"`
	OriginSecret  string `json:"sk_origin"`
	RequestBlind  string `json:"request_blind"`
	IndexRequest  string `json:"request_key"`
	IndexResponse string `json:"index_key"`
	Index         string `json:"anon_issuer_origin_id"`
}

type anonOriginIDTestVector struct {
	t             *testing.T
	curve         elliptic.Curve
	clientSecret  *ecdsa.PrivateKey
	originSecret  *ecdsa.PrivateKey
	requestBlind  *ecdsa.PrivateKey
	indexRequest  []byte
	indexResponse []byte
	index         []byte
}

type anonOriginIDTestVectorArray struct {
	t       *testing.T
	vectors []anonOriginIDTestVector
}

func (tva anonOriginIDTestVectorArray) MarshalJSON() ([]byte, error) {
	return json.Marshal(tva.vectors)
}

func (tva *anonOriginIDTestVectorArray) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &tva.vectors)
	if err != nil {
		return err
	}

	for i := range tva.vectors {
		tva.vectors[i].t = tva.t
	}
	return nil
}

func (etv anonOriginIDTestVector) MarshalJSON() ([]byte, error) {
	clientSecretKey := etv.clientSecret.D.Bytes()
	clientPublicKey := elliptic.MarshalCompressed(etv.curve, etv.clientSecret.X, etv.clientSecret.Y)
	originSecretKey := etv.originSecret.D.Bytes()
	blindKey := etv.requestBlind.D.Bytes()

	return json.Marshal(rawAnonOriginIDTestVector{
		ClientSecret:  mustHex(clientSecretKey),
		ClientPublic:  mustHex(clientPublicKey),
		OriginSecret:  mustHex(originSecretKey),
		RequestBlind:  mustHex(blindKey),
		IndexRequest:  mustHex(etv.indexRequest),
		IndexResponse: mustHex(etv.indexResponse),
		Index:         mustHex(etv.index),
	})
}

func (etv *anonOriginIDTestVector) UnmarshalJSON(data []byte) error {
	raw := rawAnonOriginIDTestVector{}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	curve := elliptic.P384()

	clientSecretKey, err := ecdsa.CreateKey(curve, mustUnhex(nil, raw.ClientSecret))
	if err != nil {
		return err
	}

	originKey, err := ecdsa.CreateKey(curve, mustUnhex(nil, raw.OriginSecret))
	if err != nil {
		return err
	}

	blindKey, err := ecdsa.CreateKey(curve, mustUnhex(nil, raw.RequestBlind))
	if err != nil {
		return err
	}

	etv.curve = curve
	etv.clientSecret = clientSecretKey
	etv.originSecret = originKey
	etv.requestBlind = blindKey
	etv.indexRequest = mustUnhex(nil, raw.IndexRequest)
	etv.indexResponse = mustUnhex(nil, raw.IndexResponse)
	etv.index = mustUnhex(nil, raw.Index)

	return nil
}

func generateAnonOriginIDTestVector(t *testing.T) anonOriginIDTestVector {
	curve := elliptic.P384()
	clientSecretKey, _ := ecdsa.GenerateKey(curve, rand.Reader)
	originSecretKey, _ := ecdsa.GenerateKey(curve, rand.Reader)
	clientBlindKey, _ := ecdsa.GenerateKey(curve, rand.Reader)

	requestKey, err := ecdsa.BlindPublicKey(curve, &clientSecretKey.PublicKey, clientBlindKey)
	if err != nil {
		t.Fatal(err)
	}

	blindedRequestKey, err := ecdsa.BlindPublicKey(curve, requestKey, originSecretKey)
	if err != nil {
		t.Fatal(err)
	}

	indexKey, err := ecdsa.UnblindPublicKey(curve, blindedRequestKey, clientBlindKey)
	if err != nil {
		t.Fatal(err)
	}

	requestKeyEnc := elliptic.MarshalCompressed(curve, requestKey.X, requestKey.Y)
	blindedRequestKeyEnc := elliptic.MarshalCompressed(curve, blindedRequestKey.X, blindedRequestKey.Y)
	clientPublicKeyEnc := elliptic.MarshalCompressed(curve, clientSecretKey.X, clientSecretKey.Y)
	indexKeyEnc := elliptic.MarshalCompressed(curve, indexKey.X, indexKey.Y)

	index, err := computeIndex(clientPublicKeyEnc, indexKeyEnc)
	if err != nil {
		t.Fatal(err)
	}

	return anonOriginIDTestVector{
		curve:         curve,
		clientSecret:  clientSecretKey,
		originSecret:  originSecretKey,
		requestBlind:  clientBlindKey,
		indexRequest:  requestKeyEnc,
		indexResponse: blindedRequestKeyEnc,
		index:         index,
	}
}

func verifyAnonOriginIDTestVector(t *testing.T, vector anonOriginIDTestVector) {
	requestKey, err := ecdsa.BlindPublicKey(vector.curve, &vector.clientSecret.PublicKey, vector.requestBlind)
	if err != nil {
		t.Fatal(err)
	}

	blindedRequestKey, err := ecdsa.BlindPublicKey(vector.curve, requestKey, vector.originSecret)
	if err != nil {
		t.Fatal(err)
	}

	indexKey, err := ecdsa.UnblindPublicKey(vector.curve, blindedRequestKey, vector.requestBlind)
	if err != nil {
		t.Fatal(err)
	}

	requestKeyEnc := elliptic.MarshalCompressed(vector.curve, requestKey.X, requestKey.Y)
	blindedRequestKeyEnc := elliptic.MarshalCompressed(vector.curve, blindedRequestKey.X, blindedRequestKey.Y)
	clientPublicKeyEnc := elliptic.MarshalCompressed(vector.curve, vector.clientSecret.X, vector.clientSecret.Y)
	indexKeyEnc := elliptic.MarshalCompressed(vector.curve, indexKey.X, indexKey.Y)

	index, err := computeIndex(clientPublicKeyEnc, indexKeyEnc)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(requestKeyEnc, vector.indexRequest) {
		t.Fatal("Index request mismatch")
	}
	if !bytes.Equal(blindedRequestKeyEnc, vector.indexResponse) {
		t.Fatal("Index response mismatch")
	}
	if !bytes.Equal(index, vector.index) {
		t.Fatal("Index mismatch")
	}
}

func verifyAnonOriginIDTestVectors(t *testing.T, encoded []byte) {
	vectors := anonOriginIDTestVectorArray{t: t}
	err := json.Unmarshal(encoded, &vectors)
	if err != nil {
		t.Fatalf("Error decoding test vector string: %v", err)
	}

	for _, vector := range vectors.vectors {
		verifyAnonOriginIDTestVector(t, vector)
	}
}

func TestVectorGenerateAnonOriginID(t *testing.T) {
	vectors := make([]anonOriginIDTestVector, 0)
	vectors = append(vectors, generateAnonOriginIDTestVector(t))

	// Encode the test vectors
	encoded, err := json.Marshal(vectors)
	if err != nil {
		t.Fatalf("Error producing test vectors: %v", err)
	}

	// Verify that we process them correctly
	verifyAnonOriginIDTestVectors(t, encoded)

	var outputFile string
	if outputFile = os.Getenv(outputAnonOriginIDTestVectorEnvironmentKey); len(outputFile) > 0 {
		err := ioutil.WriteFile(outputFile, encoded, 0644)
		if err != nil {
			t.Fatalf("Error writing test vectors: %v", err)
		}
	}
}

func TestVectorVerifyAnonOriginID(t *testing.T) {
	var inputFile string
	if inputFile = os.Getenv(inputAnonOriginIDTestVectorEnvironmentKey); len(inputFile) == 0 {
		t.Skip("Test vectors were not provided")
	}

	encoded, err := ioutil.ReadFile(inputFile)
	if err != nil {
		t.Fatalf("Failed reading test vectors: %v", err)
	}

	verifyAnonOriginIDTestVectors(t, encoded)
}
