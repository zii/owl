package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"server/base"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/vmihailenco/msgpack"
)

const (
	Minite = 60
	Hour   = 60 * Minite
	Day    = 24 * Hour
	Month  = 30 * Day
)

var ErrNil = errors.New("cache: nil returned")
var ErrNilEncoder = errors.New("cache: nil encoder")
var ErrNilDecoder = errors.New("cache: nil decoder")

func Warn(err error) {
	if err == nil {
		return
	}
	if err == ErrNil {
		return
	}
	if err == redis.ErrNil {
		return
	}
	log.Println("WARN:", err.Error())
}

var Encoders map[string]func(v interface{}) ([]byte, error)
var Decoders map[string]func(data []byte, v interface{}) error

func init() {
	Encoders = make(map[string]func(v interface{}) ([]byte, error))
	Decoders = make(map[string]func(data []byte, v interface{}) error)
	Encoders["json"] = json.Marshal
	Decoders["json"] = json.Unmarshal
	Encoders["msgpack"] = msgpack.Marshal
	Decoders["msgpack"] = msgpack.Unmarshal
}

func Encode(encoding string, v interface{}) (interface{}, error) {
	if encoding == "" {
		return v, nil
	}
	encoder, ok := Encoders[encoding]
	if !ok {
		return nil, ErrNilEncoder
	}
	return encoder(v)
}

func Decode(encoding string, b []byte, v interface{}) error {
	if encoding == "" {
		return nil
	}
	decoder, ok := Decoders[encoding]
	if !ok {
		return ErrNilDecoder
	}
	return decoder(b, v)
}

type Result struct {
	Reply    interface{}
	Err      error
	Encoding string
}

func (this *Result) Int() int {
	v, err := redis.Int(this.Reply, this.Err)
	Warn(err)
	return v
}

func (this *Result) Int32() int32 {
	v, err := redis.Int(this.Reply, this.Err)
	Warn(err)
	return int32(v)
}

func (this *Result) Int64() int64 {
	v, err := redis.Int64(this.Reply, this.Err)
	Warn(err)
	return v
}

func (this *Result) Float64() float64 {
	v, err := redis.Float64(this.Reply, this.Err)
	Warn(err)
	return v
}

func (this *Result) String() string {
	v, err := redis.String(this.Reply, this.Err)
	Warn(err)
	return v
}

func (this *Result) Bytes() []byte {
	v, err := redis.Bytes(this.Reply, this.Err)
	Warn(err)
	return v
}

func (this *Result) Unmarshal(v interface{}) error {
	b := this.Bytes()
	if len(b) <= 0 {
		return ErrNil
	}
	return Decode(this.Encoding, b, v)
}

func (this *Result) Bool() bool {
	v, err := redis.Bool(this.Reply, this.Err)
	Warn(err)
	return v
}

func (this *Result) Values() []interface{} {
	v, err := redis.Values(this.Reply, this.Err)
	Warn(err)
	return v
}

func (this *Result) Ints() []int {
	v, err := redis.Ints(this.Reply, this.Err)
	Warn(err)
	return v
}

func (this *Result) Int32s() []int32 {
	v, err := redis.Ints(this.Reply, this.Err)
	Warn(err)
	var out []int32
	for _, i := range v {
		out = append(out, int32(i))
	}
	return out
}

func (this *Result) Int64s() []int64 {
	v, err := redis.Int64s(this.Reply, this.Err)
	Warn(err)
	return v
}

func (this *Result) Float64s() []float64 {
	v, err := redis.Float64s(this.Reply, this.Err)
	Warn(err)
	return v
}

func (this *Result) Strings() []string {
	v, err := redis.Strings(this.Reply, this.Err)
	Warn(err)
	return v
}

func (this *Result) StringMap() map[string]string {
	v, err := redis.StringMap(this.Reply, this.Err)
	Warn(err)
	return v
}

func (this *Result) IntMap() map[string]int {
	v, err := redis.IntMap(this.Reply, this.Err)
	Warn(err)
	return v
}

func (this *Result) Int64Map() map[string]int64 {
	v, err := redis.Int64Map(this.Reply, this.Err)
	Warn(err)
	return v
}

func Do(dbname, command string, args ...interface{}) *Result {
	conn := GetRedisPoolClient(dbname)
	defer conn.Close()
	reply, err := conn.Do(command, args...)
	if err != nil {
		log.Println("redis-err", err, command, args)
	}
	return &Result{Reply: reply, Err: err}
}

type Key struct {
	DBName       string      // redis连接池名
	Key          interface{} // 主键名
	SecondaryKey interface{} // (可选)次级键名, 用于哈希或集合操作, 比如hget <key> <field>
	Encoding     string      // 序列化方式, msgpack/json/自定义
	Timeout      int         // 过期时间(秒)
	Critical     bool        // 是否关键key. true=只要读写报错就panic
}

func (this *Key) Do(command string, args ...interface{}) *Result {
	var newargs []interface{}
	if this.SecondaryKey != nil && strings.HasPrefix(command, "H") {
		newargs = make([]interface{}, 0, len(args)+2)
		newargs = append(newargs, this.Key, this.SecondaryKey)
		newargs = append(newargs, args...)
	} else {
		newargs = make([]interface{}, 0, len(args)+1)
		newargs = append(newargs, this.Key)
		newargs = append(newargs, args...)
	}
	result := Do(this.DBName, command, newargs...)
	if this.Critical && result.Err != nil {
		panic(result.Err)
	}
	result.Encoding = this.Encoding
	return result
}

func (this *Key) Exists() bool {
	return this.Do("EXISTS").Bool()
}

func (this *Key) Get() *Result {
	return this.Do("GET")
}

func (this *Key) Set(value interface{}) *Result {
	value, err := Encode(this.Encoding, value)
	if err != nil {
		return &Result{Reply: nil, Err: err}
	}
	if this.Timeout > 0 {
		return this.Do("SET", value, "EX", this.Timeout)
	}
	return this.Do("SET", value)
}

func (this *Key) SetEX(value interface{}, timeout int) *Result {
	value, err := Encode(this.Encoding, value)
	if err != nil {
		return &Result{Reply: nil, Err: err}
	}
	return this.Do("SET", value, "EX", timeout)
}

func (this *Key) SetNX(value interface{}) *Result {
	value, err := Encode(this.Encoding, value)
	if err != nil {
		return &Result{Reply: nil, Err: err}
	}
	if this.Timeout > 0 {
		return this.Do("SET", value, "EX", this.Timeout, "NX")
	}
	return this.Do("SET", value, "NX")
}

func (this *Key) Del() *Result {
	return this.Do("DEL")
}

func (this *Key) Incr() *Result {
	r := this.Do("INCR")
	if this.Timeout > 0 {
		this.Expire(this.Timeout)
	}
	return r
}

func (this *Key) IncrBy(amount int) *Result {
	r := this.Do("INCRBY", amount)
	if this.Timeout > 0 {
		this.Expire(this.Timeout)
	}
	return r
}

// 安全增减键值, Result.Reply==nil表示失败, Result.Int()是最新值
func (this *Key) LuaIncrBy(amount int) *Result {
	c := GetRedisPoolClient(this.DBName)
	defer c.Close()
	reply, err := lua_incrby.Do(c, this.Key, amount)
	base.Raise(err)
	return &Result{Reply: reply, Err: err}
}

func (this *Key) HExists(key ...interface{}) bool {
	return this.Do("HEXISTS", key...).Bool()
}

func (this *Key) HLen() int {
	return this.Do("HLEN").Int()
}

func (this *Key) HGet(fields ...interface{}) *Result {
	return this.Do("HGET", fields...)
}

func (this *Key) HMGet(fields ...interface{}) *Result {
	return this.Do("HMGET", fields...)
}

func (this *Key) HMSet(fields ...interface{}) *Result {
	return this.Do("HMSET", fields...)
}

func (this *Key) HGetAll() *Result {
	return this.Do("HGETALL")
}

func (this *Key) HScanEx() *Result {
	ret := &Result{}
	replys := []interface{}{}
	ret.Encoding = this.Encoding
	cursor := 0
	var err error
	for {
		ret_tmp := this.Do("HSCAN", fmt.Sprintf("%d", cursor))
		//fmt.Println(ret_tmp.Reply, "111", reflect.TypeOf(ret_tmp.Reply))
		//if ret_tmp.Err != nil {
		//	return ret_tmp
		//}
		//if ret_tmp.Reply == nil {
		//	return ret_tmp
		//}
		//if _, ok := ret_tmp.Reply.(redis.Error); ok {
		//	fmt.Println(ret_tmp.Reply, ret_tmp.Err)
		//}

		ret_objs := ret_tmp.Reply.([]interface{})
		cursor, err = redis.Int(ret_objs[0], ret_tmp.Err)
		if err != nil {
			return ret_tmp
		}
		replys = append(replys, (ret_objs[1].([]interface{}))...)
		//if ret.Reply == nil {
		//	ret.Reply = ret_objs[1].([]interface{})
		//} else {
		//	had := ret.Reply.([]interface{})
		//	ret.Reply = append(had, ret_objs[1].([]interface{}))
		//}

		if cursor == 0 {
			break
		}
	}
	ret.Reply = replys

	return ret
}

func (this *Key) HKeys(fields ...interface{}) *Result {
	return this.Do("HKEYS", fields...)
}

func (this *Key) HSet(value interface{}) *Result {
	value, err := Encode(this.Encoding, value)
	if err != nil {
		return &Result{Reply: nil, Err: err}
	}
	return this.Do("HSET", value)
}

func (this *Key) HSet2(key, value interface{}) *Result {
	value, err := Encode(this.Encoding, value)
	if err != nil {
		return &Result{Reply: nil, Err: err}
	}
	return this.Do("HSET", key, value)
}

func (this *Key) HSetNX(value interface{}) *Result {
	value, err := Encode(this.Encoding, value)
	if err != nil {
		return &Result{Reply: nil, Err: err}
	}
	return this.Do("HSETNX", value)
}

func (this *Key) HSetNX2(key, value interface{}) *Result {
	value, err := Encode(this.Encoding, value)
	if err != nil {
		return &Result{Reply: nil, Err: err}
	}
	return this.Do("HSETNX", key, value)
}

func (this *Key) HIncrBy(d int) *Result {
	r := this.Do("HINCRBY", d)
	if this.Timeout > 0 {
		this.Expire(this.Timeout)
	}
	return r
}

func (this *Key) HIncrBy2(field interface{}, d int) *Result {
	r := this.Do("HINCRBY", field, d)
	if this.Timeout > 0 {
		this.Expire(this.Timeout)
	}
	return r
}

func (this *Key) HDel(fields ...interface{}) *Result {
	return this.Do("HDEL", fields...)
}

func (this *Key) TTL() int {
	return this.Do("TTL").Int()
}

func (this *Key) Expire(ttl int) bool {
	return this.Do("EXPIRE", ttl).Bool()
}

// 有序集合添加成员, 重复添加返回false
func (this *Key) ZAdd(member, score interface{}) bool {
	return this.Do("ZADD", score, member).Bool()
}

// Only update elements that already exist.
func (this *Key) ZAddXX(member, score interface{}) bool {
	return this.Do("ZADD", "XX", score, member).Bool()
}

// Only update elements that not exist.
func (this *Key) ZAddNX(member, score interface{}) bool {
	return this.Do("ZADD", "NX", score, member).Bool()
}

// 有序集合删除成员, 重复删除返回false
func (this *Key) ZRem(member interface{}) bool {
	return this.Do("ZREM", member).Bool()
}

// Redis Zremrangebyscore 命令用于移除有序集中，指定分数（score）区间内的所有成员。
func (this *Key) ZRemRangeByScore(min interface{}, max interface{}) int {
	return this.Do("ZREMRANGEBYSCORE", min, max).Int()
}

func (this *Key) ZExists(member interface{}) bool {
	return this.Do("ZSCORE", member).Reply != nil
}

func (this *Key) ZScore(member interface{}) *Result {
	return this.Do("ZSCORE", member)
}

// d: int/float
func (this *Key) ZIncrBy(member interface{}, d interface{}) *Result {
	return this.Do("ZINCRBY", d, member)
}

func (this *Key) ZRank(member interface{}) int {
	r := this.Do("ZRANK", member)
	if r.Reply == nil {
		return -1
	}
	return r.Int()
}

func (this *Key) ZRevRank(member interface{}) int {
	r := this.Do("ZREVRANK", member)
	if r.Reply == nil {
		return -1
	}
	return r.Int()
}

func (this *Key) ZCard() int {
	r := this.Do("ZCARD")
	return r.Int()
}

// Removes and returns up to count members with the lowest scores in the sorted set stored at key.
// Returns: list of popped elements and scores.
func (this *Key) ZPopMin(count int) []string {
	r := this.Do("ZPOPMIN", count)
	return r.Strings()
}

func (this *Key) GeoDist(m1, m2 interface{}) int {
	return this.Do("GEODIST", m1, m2).Int()
}

// 返回 false 表示重复添加
func (this *Key) SAdd(members ...interface{}) bool {
	r := this.Do("SADD", members...).Bool()
	if this.Timeout > 0 {
		this.Expire(this.Timeout)
	}
	return r
}

func (this *Key) SMembers() *Result {
	return this.Do("SMEMBERS")
}

func (this *Key) SCard() *Result {
	return this.Do("SCARD")
}

// 返回 false 表示重复删除
func (this *Key) SRem(members ...interface{}) bool {
	return this.Do("SREM", members...).Bool()
}

func (this *Key) SIsMember(member interface{}) bool {
	return this.Do("SISMEMBER", member).Bool()
}

func (this *Key) SRandMember(count int) *Result {
	return this.Do("SRANDMEMBER", count)
}
func (this *Key) SPop(count int) *Result {
	return this.Do("SPOP", count)
}

func (this *Key) GeoAdd(member interface{}, lat, lng float64) bool {
	return this.Do("GEOADD", lng, lat, member).Bool()
}

// 返回(纬度, 经度)
func (this *Key) GeoPos(member interface{}) (float64, float64) {
	var lat, lng float64
	vals := this.Do("GEOPOS", member).Values()
	if len(vals) > 0 {
		p, _ := redis.Float64s(vals[0], nil)
		if len(p) == 2 {
			lat = p[1]
			lng = p[0]
		}
	}
	return lat, lng
}

// 批量返回坐标 [[纬度,经度], [纬度,经度], ...]
func (this *Key) GeoMPos(members ...interface{}) [][2]float64 {
	var out [][2]float64
	if len(members) == 0 {
		return out
	}
	vals := this.Do("GEOPOS", members...).Values()
	for _, val := range vals {
		p, _ := redis.Float64s(val, nil)
		e := [2]float64{0, 0}
		if len(p) == 2 {
			e[0], e[1] = p[0], p[1]
		}
		out = append(out, e)
	}
	return out
}

// 右进队, 返回最新队列长度
func (this *Key) RPush(args ...interface{}) int {
	var elements []interface{}
	for _, arg := range args {
		value, err := Encode(this.Encoding, arg)
		if err != nil {
			return 0
		}
		elements = append(elements, value)
	}
	return this.Do("RPUSH", elements...).Int()
}

// 左出队, 返回弹出值
func (this *Key) LPop() *Result {
	return this.Do("LPOP")
}

func (this *Key) SetBit(offset int, value int) bool {
	v := this.Do("SETBIT", offset, value).Int()
	return v != value
}

func (this *Key) BitCount() int {
	return this.Do("BITCOUNT").Int()
}
