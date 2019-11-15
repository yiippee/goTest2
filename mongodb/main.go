// mgotest project main.go
package main

import (
	"fmt"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type User struct {
	Id        bson.ObjectId `bson:"_id"`
	Username  string        `bson:"name"`
	Pass      string        `bson:"pass"`
	Regtime   int64         `bson:"regtime"`
	Interests []string      `bson:"interests"`
}

type Student struct {
	Id    bson.ObjectId `bson:"_id"`
	Sname string        `bson:"sname"`
	Sage  int           `bson:"sage"`
	Score float64       `bson:"score"`
	Records []Record    `bson:"records"`
	Time  time.Time     `bson:"time"`
}
type Record struct {
	A int `bson:"a"`
	B int `bson:"b"`
	C string `bson:"c"`
}

type WriteConcern struct {
	W int `bson:"w"`
}

const URL string = "127.0.0.1:27017"

var c *mgo.Collection
var session *mgo.Session

func (user User) ToString() string {
	return fmt.Sprintf("%#v", user)
}

func init() {
	session, _ = mgo.Dial(URL)
	//切换到数据库
	db := session.DB("blog")
	//切换到collection
	c = db.C("student")
}
func myFind() {
	//    defer session.Close()
	var stu []Student
	c.Find(bson.M{"sname": "lisi"}).All(&stu)
	//for _, value := range stu {
	//	fmt.Println(value)
	//}
	////根据ObjectId进行查询
	//idStr := stu[0].Id
	//objectId := idStr
	//user := new(Student)
	//c.Find(bson.M{"_id": objectId}).One(user)
	//fmt.Println(user)
}

func myInsert() {
	stu := Student{
		Id: bson.NewObjectId(),
		Sname: "lizhanbin",
		Sage:  31,
		Score: 100,
		Time: time.Now(),
		Records: []Record{
			{1,2,"12"},
			{3, 4, "34"},
			{5, 6, "56"},
		},
	}
	err := c.Insert(stu)
	if err == nil {
		fmt.Println("插入成功")
	} else {
		fmt.Println(err.Error())
		defer panic(err)
	}

	// c.Upsert()
}

//新增数据
func add() {
	//    defer session.Close()
	stu1 := new(User)
	stu1.Id = bson.NewObjectId()
	stu1.Username = "stu1_name"
	stu1.Pass = "stu1_pass"
	stu1.Regtime = time.Now().Unix()
	stu1.Interests = []string{"象棋", "游泳", "跑步"}
	err := c.Insert(stu1)
	if err == nil {
		fmt.Println("插入成功")
	} else {
		fmt.Println(err.Error())
		defer panic(err)
	}
}

//查询
func find() {
	//    defer session.Close()
	var users []User
	//    c.Find(nil).All(&users)
	c.Find(bson.M{"name": "stu1_name"}).All(&users)
	for _, value := range users {
		fmt.Println(value.ToString())
	}
	//根据ObjectId进行查询
	idStr := "577fb2d1cde67307e819133d"
	objectId := bson.ObjectIdHex(idStr)
	user := new(User)
	c.Find(bson.M{"_id": objectId}).One(user)
	fmt.Println(user)
}

//根据id进行修改
func update() {
	interests := []string{"象棋", "游泳", "跑步"}
	err := c.Update(bson.M{"_id": bson.ObjectIdHex("577fb2d1cde67307e819133d")}, bson.M{"$set": bson.M{
		"name":      "修改后的name",
		"pass":      "修改后的pass",
		"regtime":   time.Now().Unix(),
		"interests": interests,
	}})
	if err != nil {
		fmt.Println("修改失败")
	} else {
		fmt.Println("修改成功")
	}
}

//删除
func del() {
	err := c.Remove(bson.M{"_id": bson.ObjectIdHex("577fb2d1cde67307e819133d")})
	if err != nil {
		fmt.Println("删除失败" + err.Error())
	} else {
		fmt.Println("删除成功")
	}
}
func main() {
	myFind()
	myInsert()
	return
	add()
	find()
	update()
	del()
}
