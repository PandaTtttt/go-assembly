package simplejson_test

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/PandaTtttt/go-assembly/env/conn"
	"github.com/PandaTtttt/go-assembly/env/mysql"
	"github.com/PandaTtttt/go-assembly/simplejson"
	"github.com/PandaTtttt/go-assembly/util/must"
	"testing"

	"gorm.io/gorm"
	. "gorm.io/gorm/utils/tests"
)

var _ driver.Valuer = &simplejson.JSON{}

func TestJSON(t *testing.T) {
	mysql.Init(&mysql.Params{
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "hello",
		User:     "root",
		Password: "1024",
	})

	type UserWithJSON struct {
		gorm.Model
		Name       string
		Attributes *simplejson.JSON
	}

	must.Must(conn.DB().Migrator().DropTable(&UserWithJSON{}))

	if err := conn.DB().Migrator().AutoMigrate(&UserWithJSON{}); err != nil {
		t.Errorf("failed to migrate, got error: %v", err)
	}

	user1Attrs := `{"age":18,"name":"json-1","orgs":{"orga":"orga"},"tags":["tag1","tag2"]}`
	user2Attrs := `{"name": "json-2", "age": 28, "tags": ["tag1", "tag3"], "role": "admin", "orgs": {"orgb": "orgb"}}`

	user1Json, err := simplejson.NewJSON([]byte(user1Attrs))
	must.Must(err)
	user2Json, err := simplejson.NewJSON([]byte(user2Attrs))
	must.Must(err)

	users := []UserWithJSON{{
		Name:       "json-1",
		Attributes: user1Json,
	}, {
		Name:       "json-2",
		Attributes: user2Json,
	}}

	if err := conn.DB().Create(&users).Error; err != nil {
		t.Errorf("Failed to create users %v", err)
	}

	var result UserWithJSON
	if err := conn.DB().First(&result, simplejson.JSONQuery("attributes").HasKey("role")).Error; err != nil {
		t.Fatalf("failed to find user with json key, got error %v", err)
	}
	AssertEqual(t, result.Name, users[1].Name)

	var result2 UserWithJSON
	if err := conn.DB().Where(simplejson.JSONQuery("attributes").HasKey("orgs", "orga")).Find(&result2).Error; err != nil {
		t.Fatalf("failed to find user with json key, got error %v", err)
	}
	AssertEqual(t, result2.Name, users[0].Name)

	result2Attrs, err := json.Marshal(&result2.Attributes)
	if err != nil {
		t.Fatalf("failed to marshal result2.Attributes, got error %v", err)
	}
	AssertEqual(t, string(result2Attrs), user1Attrs)

	var j simplejson.JSON
	if err := j.UnmarshalJSON([]byte(user1Attrs)); err != nil {
		t.Fatalf("failed to unmarshal user1Attrs, got error %v", err)
	}
	b, err := j.MarshalJSON()
	must.Must(err)

	AssertEqual(t, string(b), user1Attrs)

	var result3 UserWithJSON
	if err := conn.DB().First(&result3, simplejson.JSONQuery("attributes").Equals("json-1", "name")).Error; err != nil {
		t.Fatalf("failed to find user with json value, got error %v", err)
	}
	AssertEqual(t, result3.Name, users[0].Name)

	var result4 UserWithJSON
	if err := conn.DB().Where(simplejson.JSONQuery("attributes").Equals("orgb", "orgs", "orgb")).Find(&result4).Error; err != nil {
		t.Fatalf("failed to find user with json value, got error %v", err)
	}
	AssertEqual(t, result4.Name, users[1].Name)
}
