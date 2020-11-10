package main

import (
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Chat [...]
type Chat struct {
	ID        int64     `gorm:"primary_key;column:id;type:bigint(20);not null"`
	Sender    int64     `gorm:"column:sender;type:bigint(20);not null"`
	Receiver  int64     `gorm:"column:receiver;type:bigint(20);not null"`
	Status    int       `gorm:"column:status;type:int(11);not null"` // 0:无效；1：有效
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;not null"`
}

// Events [...]
type Events struct {
	ID         int64     `gorm:"primary_key;column:id;type:bigint(20);not null"`
	UserID     int       `gorm:"column:user_id;type:int(11) unsigned"` // 事件执行者
	Operation  string    `gorm:"column:operation;type:varchar(255)"`
	TableName  string    `gorm:"column:table_name;type:varchar(255)"`
	TableID    int64     `gorm:"column:table_id;type:bigint(20)"`
	ColumnName string    `gorm:"column:column_name;type:varchar(255)"`
	OldState   string    `gorm:"column:old_state;type:varchar(255)"`
	NewState   string    `gorm:"column:new_state;type:varchar(255)"`
	Type       string    `gorm:"column:type;type:varchar(255)"` // 它属于哪个子类?e.g. CreateVideoFeedEvent, DeleteVideoFeedEvent, UpdateWikiFeedEvent
	IsDeleted  bool      `gorm:"column:is_deleted;type:tinyint(1)"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp"`
}

// FollowFeeds [...]
type FollowFeeds struct {
	ID       int64     `gorm:"primary_key;column:id;type:bigint(20);not null"`
	UserID   int       `gorm:"column:user_id;type:int(11);not null"`
	EventID  int       `gorm:"column:event_id;type:int(11);not null"`
	RoutedAt time.Time `gorm:"column:routed_at;type:timestamp;not null"`
}

// Friend [...]
type Friend struct {
	ID        int64     `gorm:"primary_key;column:id;type:bigint(20);not null"`
	UserID1   int64     `gorm:"column:user_id_1;type:bigint(20);not null"`
	UserID2   int64     `gorm:"column:user_id_2;type:bigint(20);not null"`
	Status    int       `gorm:"column:status;type:int(11);not null"` // 0:没关系；1：小关注大；2：大关注小
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;not null"`
}

// Likes [...]
type Likes struct {
	ID        int       `gorm:"primary_key;column:id;type:int(11);not null"`
	IP        string    `gorm:"column:ip;type:varchar(20);not null"`
	Name      string    `gorm:"column:name;type:varchar(256);not null"`
	Title     string    `gorm:"column:title;type:varchar(128);not null"`
	Hash      int64     `gorm:"column:hash;type:bigint(20) unsigned"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;not null"`
}

// Message [...]
type Message struct {
	ID        int64     `gorm:"primary_key;column:id;type:bigint(20);not null"`
	Msg       string    `gorm:"column:msg;type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:timestamp;not null"`
}

// ObjectFeeds [...]
type ObjectFeeds struct {
	ID         int64     `gorm:"primary_key;column:id;type:bigint(20);not null"`
	ObjectID   int       `gorm:"column:object_id;type:int(11);not null"`
	ObjectType string    `gorm:"column:object_type;type:varchar(255)"`
	EventID    int       `gorm:"column:event_id;type:int(11);not null"`
	RoutedAt   time.Time `gorm:"column:routed_at;type:datetime;not null"`
}

// TagFeeds tag_feeds 来储存「与某个 tag 有关的 events」。
type TagFeeds struct {
	ID       int64     `gorm:"primary_key;column:id;type:bigint(20);not null"`
	TagID    int       `gorm:"column:tag_id;type:int(11);not null"`
	EventID  int       `gorm:"column:event_id;type:int(11);not null"`
	RoutedAt time.Time `gorm:"column:routed_at;type:datetime;not null"`
}

var db *gorm.DB

func init() {
	//创建 gorm.DB
	var err error
	db, err = gorm.Open("mysql", "dev:ADq3rfaqhwer6#!a76f!@tcp(192.168.216.46:3306)/lzb-test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
}
func main() {
	/*
			三张表来描述用户信息
			user
			user_auth
			user_detail   记录一下最后一次im的时间，这样就能按照时间来查询chat记录了
		                  记录一下最后发送新鲜事的周数，便于查表

	*/

	// 创建好友关系  三个好友： 111   222   333  互相都为好友
	// INSERT INTO friend (user_id_1, user_id_2, STATUS) values (222, 333, 2)  偏序

	// 111 发送一条消息给 222  内容为 "你好哇，今天天气很好。"

	/*
		insert into message_0001 (msg) values ("你好哇，吃饭了么?hahaha");    msg表可按照500万一张分表，需要有一个postion表，记录当前分表到哪了。
		insert into chat_202011 (sender, receiver, msg_id, status) values (111, 222, 2, 1);
		chat 表可按照 sender+receiver 分表，两者唯一确定一条聊天记录。     可以按时间来分表。 这个时间可以加在user表中。

		其实聊天记录是很少去服务端查询的，可以app本地存储。
	*/

	/*

				    feed 流设计

					select * from follow_feeds where user_id = 111 order by routed_at desc limit 0, 3;
				按时间排序，从用户 111 的收件箱中查询3条最新的信息


				follow_feeds_111 : xxx,  yyy,  zzz,  mmm,  nnn   链表

				follow_feeds
			           111  1   xxx
		               111  2   yyy
		               111  3   zzz
		               111  3   mmm
		               111  3   nnn

	*/
}
