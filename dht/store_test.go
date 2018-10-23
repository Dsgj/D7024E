package dht

import (
	"fmt"
	"testing"
	"time"
	"math/rand"
	//"encoding/hex"
)

func TestStore(t *testing.T) {
	
	contacts, _, teardown := setupTestCase(t, 20)
	
	defer teardown(t)

	rand.Seed(time.Now().UnixNano())
	N := rand.Intn(9)
	N += 1
	data := make([]byte, N)
	for i := 0; i < N; i++ {
		data[i] = 'a' + byte(i%26)
	}
	
	store := NewStore()
	
	rec := store.Store(data, contacts[1], time.Now())
	
	for i := 0; i < len(data); i++ {
		if data[i] != rec.value[i] {
			t.Errorf("Store stored the wrong value. Stored: %d. Should store: %d.\n", rec.value, data)
		}
	}
	if rec.publisher != contacts[1] {
		t.Errorf("Store stored the wrong publisher. Stored %d. Should store: %d.\n", rec.publisher, contacts[1])
	}
	
}

func Testfromstring(t *testing.T) {
	//Tests FromString(), the input has to be long enough so that the for-loop doesnt get a outofbounds error with the bytearray[20]
	fmt.Println("len > 20")
	FromString("48656c6c6f20476f706865722148656c6c6f20476f7068657221")
	fmt.Println(FromString("48656c6c6f20476f706865722148656c6c6f20476f7068657221"))
	fmt.Println("\nlen == 20")
	FromString("48656c6c6f20476f706865722148656c6c6f20476")
	fmt.Println(FromString("48656c6c6f20476f706865722148656c6c6f20476"))
	fmt.Println("\nlen < 20")
	FromString("48656c6c6f20476f7068657221")
	fmt.Println(FromString("48656c6c6f20476f7068657221"))
}
