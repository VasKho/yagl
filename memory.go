package yagl

import (
	md "crypto/md5"
	"fmt"
	"sync"
)

const (
	// Element types
	Node_t int = 0x1
	Arc_t  int = 0x2
	// Element constancy
	Const_t int = 0x4
	Var_t   int = 0x8
	// Arc types
	Pos_t  int = 0x10
	Neg_t  int = 0x20
	Fuz_t  int = 0x40
	Temp_t int = 0x80
	Perm_t int = 0x100
	// Node types
	Tuple_t  int = 0x10
	Struct_t int = 0x20
	Role_t   int = 0x40
	Norole_t int = 0x80
	Class_t  int = 0x100
)

type Exchange[T any] struct {
	Result chan T
	Wg     *sync.WaitGroup
}

func genHash(identifier string, el_type int) string {
	bytes := []byte(identifier)
	bytes = append(bytes, byte(el_type))
	hash_value := md.Sum(bytes)
	return fmt.Sprintf("%x", hash_value)
}
