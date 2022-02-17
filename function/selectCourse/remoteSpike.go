package selectCourse

import (
	"github.com/garyburd/redigo/redis"
)

const LuaScript = `
        local course_id = KEYS[1]
        local user_id = ARGV[1]
		-- 课程库存key
		local product_stock_key = 'seckill:' .. course_id .. ':stock'
		-- 课程秒杀结束标识的key
		local end_product_key = 'seckill:' .. course_id .. ':end'

		-- 存储秒杀成功的用户id的集合的key
		local bought_users_key = 'seckill:' .. course_id .. ':uids'
		local students_list = 'studentsID'
		-- 存储用户抢到的课
		local courses_key = 'courses:' .. user_id .. ':uids'
		-- 判断学生是否存在
		local is_student = redis.call('bf.exists',students_list,user_id)
		if is_student <= 0 then
			return -3
		end
		--判断该商品是否秒杀结束
		local is_end = redis.call('get',end_product_key)
		-- 判断课程是否存在
		if not is_end then
			return -1
		end
		if  is_end and tonumber(is_end) == 1 then
    		return -2
		end
		-- 判断用户是否秒杀过
		local is_in = redis.call('bf.exists',bought_users_key,user_id)

		if is_in > 0 then
    		return 0
		end

		-- 获取商品当前库存
		local stock = redis.call('get',product_stock_key)

		-- 如果库存<=0,则返回-1
		if not stock or tonumber(stock) <=0 then
    		redis.call("set",end_product_key,"1")
			redis.call("del",bought_users_key)
    		return -1
		end

		-- 减库存,并且把用户的id添加进已购买用户set里,并给用户添加课程
		redis.call("decr",product_stock_key)
		redis.call("bf.add",bought_users_key,user_id)
		redis.call("sadd",courses_key,course_id)
		return 1
`

//初始化redis连接池
func NewPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   12000,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", ":6379")
			if err != nil {
				panic(err.Error())
			}
			//_, err = c.Do("AUTH", "root")
			_, err = c.Do("AUTH", "bytedancecamp")
			if err != nil {
				c.Close()
				panic(err.Error())
			}
			return c, err
		},
	}
}

//远端统一扣库存
func RemoteDeductionStock(conn redis.Conn, cid string, uid string) int {
	lua := redis.NewScript(1, LuaScript)
	result, err := redis.Int(lua.Do(conn, cid, uid))
	//fmt.Println(result)
	if err != nil {
		panic(err.Error())
		return -2
	}
	return result
}
