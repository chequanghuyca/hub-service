package common

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/btcsuite/btcutil/base58"
)

type UID struct {
	localID    uint32
	objectType int
	shardID    uint32
}

func NewUID(
	localID uint32,
	objectType int,
	shardID uint32,
) UID {
	return UID{
		localID:    localID,
		objectType: objectType,
		shardID:    shardID,
	}
}

// SET BITS
// Shard: 1, ObjectType: 1, ID: 1  ==> 0001 0001 0001
// 1 << 8 = 0001 0000 0000
// 1 << 4 = 0000 0001 0000
// 1 << 0 = 0000 0000 0001
func (uid UID) String() string {
	val := uint64(uid.localID)<<28 | uint64(uid.objectType)<<18 | uint64(uid.shardID)<<0
	return base58.Encode([]byte(fmt.Sprintf("%v", val)))
}

func (uid UID) GetLocalID() uint32 {
	return uid.localID
}

func (uid UID) GetShardID() uint32 {
	return uid.shardID
}

func (uid UID) GetObjectType() int {
	return uid.objectType
}

func DecomposeUID(s string) (UID, error) {
	uid, err := strconv.ParseUint(s, 10, 64)

	if err != nil {
		return UID{}, err
	}

	if (1 << 18) > uid {
		return UID{}, errors.New("invalid uid")
	}

	u := UID{
		localID:    uint32(uid >> 28),
		objectType: int(uid >> 18 & 0x3FF),
		shardID:    uint32(uid >> 0 & 0x3FFFF),
	}

	return u, nil
}

func FromBase58(s string) (UID, error) {
	return DecomposeUID(string(base58.Decode(s)))
}

func (uid UID) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", uid.String())), nil
}

func (uid *UID) UnmarshalJSON(data []byte) error {
	decodeUID, err := FromBase58(strings.Replace(string(data), "\"", "", -1))

	if err != nil {
		return err
	}

	uid.localID = decodeUID.localID
	uid.objectType = decodeUID.objectType
	uid.shardID = decodeUID.shardID

	return nil
}

func (uid *UID) Value() (driver.Value, error) {
	if uid == nil {
		return nil, nil
	}

	return int64(uid.localID), nil
}

func (uid *UID) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var i uint32

	switch v := value.(type) {
	case int:
		i = uint32(v)
	case int8:
		i = uint32(v)
	case int16:
		i = uint32(v)
	case int32:
		i = uint32(v)
	case int64:
		i = uint32(v)
	case []byte:
		a, err := strconv.Atoi(string(v))
		if err != nil {
			return err
		}

		i = uint32(a)
	default:
		return fmt.Errorf("incompatible type for UID")
	}

	*uid = NewUID(i, 0, 0)

	return nil
}
