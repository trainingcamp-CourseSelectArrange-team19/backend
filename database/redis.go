package database

import (
	"context"
	"errors"
	"strings"
	"techtrainingcamp-courseSelectArrange/tools"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()
var rdb *redis.Client

const EPTIME = 30

func RedisInitClient() {
	//初始化客户端
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: 40,
	})
	// ctx = context.Background()
}

func RedisClose() {
	rdb.Close()
}

func RedisUpdateDownloadStatus(ruleid string, status bool) error {
	//RedisInitClient()
	//defer rdb.Close()
	hits, err := rdb.HGet(ctx, ruleid, "hit_count").Result()
	hits = ToStr(ToInt(hits) + 1)
	err = rdb.HSet(ctx, ruleid, "hit_count", hits).Err()
	if status {
		downs, _ := rdb.HGet(ctx, ruleid, "download_count").Result()
		downs = ToStr(ToInt(downs) + 1)
		rdb.HSet(ctx, ruleid, "download_count", downs)
	}
	rdb.Expire(ctx, ruleid, EPTIME*time.Second)
        rdb.Expire(ctx, ruleid+"s", EPTIME*time.Second)
	return err
}

func RedisQueryRuleByID(ruleid string) (*[]map[string]string, *[]string, error) {
	//RedisInitClient()
	//defer rdb.Close()
	val, err := rdb.HGetAll(ctx, ruleid).Result()
	//pipe.Expire(ctx, ruleid, EPTIME*time.Second)
	//pipe.Expire(ctx, ruleid+"s", EPTIME*time.Second)
	// val, _ := res[0].(*redis.StringStringMapCmd).Result()
	devices := make([]map[string]string, 0)
	s := strings.Split(val["device_list"], ",")
	if len(val) == 0 {
		err = errors.New("Can't find in redis...")
		return &devices, &s, err
	}
	checkErr(err)
	devices = append(devices, val)
	return &devices, &s, err
}

func RedisDeleteRule(ruleid string) error {
	//RedisInitClient()
	//defer rdb.Close()
	pipe := rdb.TxPipeline()
	pipe.SRem(ctx, "IDList", ruleid).Err()
	pipe.Del(ctx, ruleid).Err()
	pipe.Del(ctx, ruleid+"s").Err()
	_, err := pipe.Exec(ctx)
	checkErr(err)
	return err
}

func RedisTouchRule(ruleid string) {
	//RedisInitClient()
	//defer rdb.Close()
	err := rdb.SAdd(ctx, "IDList", ruleid).Err()
	if err != nil {
		tools.LogMsg(err)
		panic(err)
	}
}

//Redis 更新规则，如果没有则创建，有则覆盖
func RedisUpdateRule(ruleid string, r *map[string]string, devices *[]string) error {
	//RedisInitClient()
	//defer rdb.Close()
	pipe := rdb.TxPipeline()
	err := pipe.SAdd(ctx, "IDList", ruleid).Err()
	checkErr(err)
	err = pipe.HMSet(ctx, ruleid, *r).Err()
	checkErr(err)
	pipe.Expire(ctx, ruleid, EPTIME*time.Second)
	//s := strings.Split(r["device_list"], ",")
	if devices != nil {
		pipe.Del(ctx, ruleid+"s")
		err = pipe.SAdd(ctx, ruleid+"s", *devices).Err()
		checkErr(err)
		pipe.Expire(ctx, ruleid+"s", EPTIME*time.Second)
	}
	_, err = pipe.Exec(ctx)
	return err
}

func RedisUpdateRuleWithList(ruleid string, r *map[string]string) error {
	s := strings.Split((*r)["device_list"], ",")
	return RedisUpdateRule(ruleid, r, &s)
}

func RedisGetRuleAttr(ruleid string, attrcode string) (string, error) {
	//RedisInitClient()
	//defer rdb.Close()
	//pipe := rdb.TxPipeline()
	val, err := rdb.HGet(ctx, ruleid, attrcode).Result()
	//pipe.Expire(ctx, ruleid, EPTIME*time.Second)
	// res, _ := pipe.Exec(ctx)
	// val, err := res[0].(*redis.StringCmd).Result()
	return val, err

}

func RedisCheckWhiteList(ruleid string, userid string) (bool, error) {
	//RedisInitClient()
	//defer rdb.Close()
	//pipe := rdb.TxPipeline()
	val, err := rdb.SIsMember(ctx, ruleid+"s", userid).Result()
	res, err := rdb.Exists(ctx, ruleid+"s").Result()
	if res == 0 {
		err = errors.New("ruleid doesn't exist!")
	}
	//pipe.Expire(ctx, ruleid+"s", EPTIME*time.Second)
	//pipe.Expire(ctx, ruleid, EPTIME*time.Second)
	return val, err
}

func GetIDList() (*[]string, error) {
	//RedisInitClient()
	//defer rdb.Close()
	val, err := rdb.SMembers(ctx, "IDList").Result()
	checkErr(err)
	return &val, err
}

func RedisDeleteAll() {
	//RedisInitClient()
	//defer rdb.Close()
	rdb.FlushAll(ctx)
}

func RedisGetAllKeys() []string {
	//RedisInitClient()
	//defer rdb.Close()
	str, _ := rdb.Keys(ctx, "*").Result()
	return str
}

// func RedisAddRule(r map[string]string, white_list []string) error {
// 	err := rdb.HMSet(ctx, strconv.Itoa(cur_id), r).Err()
// 	if err != nil {
// 		return err
// 	}
// 	err = rdb.SAdd(ctx, strconv.Itoa(cur_id)+"s", white_list).Err()
// 	if err != nil {
// 		return err
// 	}
// 	cur_id++
// 	rdb.Incr(ctx, "cur_id")
// 	return err
// }
