package handlers

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	sdk "github.com/binance-chain/go-sdk/common/types"
	"github.com/cosmostation/mintscan-binance-dex-backend/mintscan/client"
	"github.com/cosmostation/mintscan-binance-dex-backend/mintscan/db"
	"github.com/cosmostation/mintscan-binance-dex-backend/mintscan/errors"
	"github.com/cosmostation/mintscan-binance-dex-backend/mintscan/schema"
	"github.com/cosmostation/mintscan-binance-dex-backend/mintscan/utils"
	ethcmn "github.com/ethereum/go-ethereum/common"
)

// TwoAuth is a TwoAuth handler
type TwoAuth struct {
	l      *log.Logger
	client *client.Client
	db     *db.Database
}

// NewTwoAuth creates a new TwoAuth handler with the given params
func NewTwoAuth(l *log.Logger, client *client.Client, db *db.Database) *TwoAuth {
	return &TwoAuth{l, client, db}
}

// GetTwoAuth creae new key and save DB, then return TwoAuth on the active chain
func (ta *TwoAuth) Auth(rw http.ResponseWriter, r *http.Request) {
	var id, passwd int

	if len(r.URL.Query()["id"]) > 0 {
		id, _ = strconv.Atoi(r.URL.Query()["id"][0])
	} else {
		ta.l.Printf("failed to get id")
		return
	}
	if len(r.URL.Query()["passwd"]) > 0 {
		passwd, _ = strconv.Atoi(r.URL.Query()["passwd"][0])
	} else {
		ta.l.Printf("failed to get passwd")
		return
	}

	twoAuthInfo, err := ta.db.QueryTwoAuthByID(int64(id))
	if err != nil {
		ta.l.Printf("failed to query twoAuthInfo: %s", err)
		return
	}
	//fmt.Println(id, passwd, twoAuthInfo.Key)

	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(twoAuthInfo.Key)
	epochSeconds := time.Now().Unix()
	epochSeconds /= 30
	//fmt.Println(epochSeconds)
	indexs := []int64{0, 1, -1, 2, -2}
	hit := false
	for i := 0; i < len(indexs); i++ {
		epoch := toBytes(epochSeconds + indexs[i])
		pwd := oneTimePassword(key, epoch)
		//fmt.Println(epochSeconds+indexs[i], pwd)
		if pwd == uint32(passwd) {
			hit = true
			break
		}
	}
	//fmt.Println(hit)

	utils.Respond(rw, hit)
	return
}

// GetTwoAuth creae new key and save DB, then return TwoAuth on the active chain
func (ta *TwoAuth) AuthForDCI(rw http.ResponseWriter, r *http.Request) {
	var address string
	var passwd int

	if len(r.URL.Query()["address"]) > 0 {
		address = r.URL.Query()["address"][0]
	} else {
		ta.l.Printf("failed to get id")
		return
	}
	if len(r.URL.Query()["passwd"]) > 0 {
		passwd, _ = strconv.Atoi(r.URL.Query()["passwd"][0])
	} else {
		ta.l.Printf("failed to get passwd")
		return
	}

	twoAuthInfo, err := ta.db.QueryTwoAuthByAddressForDCI(address)
	if err != nil {
		ta.l.Printf("failed to query twoAuthInfo: %s", err)
		return
	}
	//fmt.Println(id, passwd, twoAuthInfo.Key)

	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(twoAuthInfo.Key)
	epochSeconds := time.Now().Unix()
	epochSeconds /= 30
	//fmt.Println(epochSeconds)
	indexs := []int64{0, 1, -1, 2, -2}
	hit := false
	for i := 0; i < len(indexs); i++ {
		epoch := toBytes(epochSeconds + indexs[i])
		pwd := oneTimePassword(key, epoch)
		//fmt.Println(epochSeconds+indexs[i], pwd)
		if pwd == uint32(passwd) {
			hit = true
			if !twoAuthInfo.Bind {
				twoAuthInfo.Bind = true
				err = ta.db.UpdateTwoAuthForDCI(twoAuthInfo)
				if err != nil {
					ta.l.Printf("failed to update TwoAuth : %s", err)
					return
				}
			}

			break
		}
	}
	//fmt.Println(hit)

	utils.Respond(rw, hit)
	return
}

func toBytes(value int64) []byte {
	var result []byte
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	for _, shift := range shifts {
		result = append(result, byte((value>>shift)&mask))
	}
	return result
}

func toUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
		(uint32(bytes[2]) << 8) + uint32(bytes[3])
}

func oneTimePassword(key []byte, value []byte) uint32 {
	// sign the value using HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, key)
	hmacSha1.Write(value)
	hash := hmacSha1.Sum(nil)

	// We're going to use a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	offset := hash[len(hash)-1] & 0x0F

	// get a 32-bit (4-byte) chunk from the hash starting at offset
	hashParts := hash[offset : offset+4]

	// ignore the most significant bit as per RFC 4226
	hashParts[0] = hashParts[0] & 0x7F

	number := toUint32(hashParts)

	// size to 6 digits
	// one million is the first number with 7 digits so the remainder
	// of the division will always return < 7 digits
	pwd := number % 1000000

	return pwd
}

var (
	Table = []string{
		"A", "B", "C", "D", "E", "F", "G", "H", // 7
		"I", "J", "K", "L", "M", "N", "O", "P", // 15
		"Q", "R", "S", "T", "U", "V", "W", "X", // 23
		"Y", "Z", "2", "3", "4", "5", "6", "7", // 31
		// padding char
	}
)

func CreateSecret() string {
	var (
		length int = 32
		secret []string
	)

	len := big.NewInt(int64(len(Table)))
	for i := 0; i < length; i++ {
		randNum, _ := rand.Int(rand.Reader, len)
		secret = append(secret, Table[randNum.Int64()])
	}
	return strings.Join(secret, "")
}

// Generate returns TwoAuth information
func (ta *TwoAuth) GenerateForDCI(rw http.ResponseWriter, r *http.Request) {

	var address string

	if len(r.URL.Query()["address"]) > 0 {
		address = r.URL.Query()["address"][0]
	} else {
		ta.l.Printf("failed to get address")
		return
	}

	bind, existed, err := ta.db.QueryTwoAuthExistedByAddress(address)
	if err != nil {
		ta.l.Printf("query address err: %s", err)
		return
	}

	if bind {
		utils.Respond(rw, struct {
			Bind bool `json:bind`
		}{
			Bind: true,
		})
		return
	}

	twoAuthInfo := schema.TwoAuthForDCI{
		Key:     CreateSecret(),
		Address: address,
	}

	if existed {
		//utils.Respond(rw, true)
		err = ta.db.UpdateTwoAuthForDCI(&twoAuthInfo)
		if err != nil {
			ta.l.Printf("failed to update TwoAuth : %s", err)
			return
		}
	} else {

		err = ta.db.InsertTwoAuthForDCI(&twoAuthInfo)
		if err != nil {
			ta.l.Printf("failed to insert TwoAuth : %s", err)
			return
		}
	}
	//fmt.Println(twoAuthInfo.ID, twoAuthInfo.Key)
	utils.Respond(rw, twoAuthInfo)
	return
}

// Generate returns TwoAuth information
func (ta *TwoAuth) Generate(rw http.ResponseWriter, r *http.Request) {
	twoAuthInfo := schema.TwoAuth{Key: CreateSecret()}

	err := ta.db.InsertTwoAuth(&twoAuthInfo)
	if err != nil {
		ta.l.Printf("failed to insert TwoAuth : %s", err)
		return
	}
	//fmt.Println(twoAuthInfo.ID, twoAuthInfo.Key)
	utils.Respond(rw, twoAuthInfo)
	return
}

func (ta *TwoAuth) GetHalleByEth(rw http.ResponseWriter, r *http.Request) {

	var q_address string
	if len(r.URL.Query()["address"]) > 0 {
		q_address = r.URL.Query()["address"][0]
	} else {
		//address 为必填项
		errors.ErrInvalidFormat(rw, http.StatusBadRequest)
		return
	}

	if len(q_address) == 42 {
		q_address = q_address[2:]
	}

	acc, err := sdk.AccAddressFromHex(q_address)
	if err != nil {
		ta.l.Printf("failed to AccAddressFromHex : %s", err)
		return
	}

	utils.Respond(rw, acc)
}

func (ta *TwoAuth) GetEthByHalle(rw http.ResponseWriter, r *http.Request) {

	var q_address string
	if len(r.URL.Query()["address"]) > 0 {
		q_address = r.URL.Query()["address"][0]
	} else {
		//address 为必填项
		errors.ErrInvalidFormat(rw, http.StatusBadRequest)
		return
	}

	addr, err := sdk.AccAddressFromBech32(q_address)
	if err != nil {
		ta.l.Printf("failed to AccAddressFromBech32 : %s", err)
		return
	}
	address := ethcmn.BytesToAddress(addr.Bytes())
	utils.Respond(rw, address)
}
