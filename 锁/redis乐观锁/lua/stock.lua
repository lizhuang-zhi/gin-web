local key = KEYS[1]   -- redis的key（lua数组索引从1开始）
local quantity = tonumber(ARGV[1])  -- 5（要扣减的数量）
local expire = tonumber(ARGV[2])  -- 10（过期时间，秒）

local current = tonumber(redis.call('GET', key) or "0")  -- 获取当前库存逻辑

if current < quantity then
    return {0, current}  -- 第一个元素表示状态码（在 Go 中会转为 []interface{}）
end

redis.call('DECRBY', key, quantity)  -- DECRBY 是 Redis 原子操作命令

if expire > 0 then
    redis.call('EXPIRE', key, expire)  -- 设置key的过期时间
end

return {1, current - quantity}  -- 1表示成功，0表示失败