-- 获取当前时间戳
local now = tonumber(redis.call("TIME")[1])
-- 如果时区是 UTC+8，那么 UTC+0 时间戳转换成UTC+8，就需要加上 8 小时的偏移量，即 (now - (now % 86400)) + 8*3600。这样获取的时间戳就是UTC+8 0 点整的时间戳了。
-- 获取今天零点时间戳如果不转换就是八点
local timezone_offset = 8 * 3600 -- 时区偏移量，单位是秒
local today = now - (now % 86400)

-- 获取已发送次数和最后发送时间
local count, lastTime, expireTime = unpack(redis.call("HMGET", KEYS[1], "count", "last_time"))

if not count then
    -- 如果没有记录则初始化
    count = 0
    lastTime = 0
    expireTime = tonumber(ARGV[3])
    redis.call("HMSET", KEYS[1], "count", count, "last_time", lastTime)
    redis.call("EXPIREAT", KEYS[1], expireTime)
else
    count = tonumber(count)
    lastTime = tonumber(lastTime)
    expireTime = tonumber(expireTime)
end

local interval = tonumber(ARGV[1]) -- 间隔时间
local limit = tonumber(ARGV[2]) -- 每天限制次数

-- 计算距离下次发送需要等待的时间
local waitTime = 0
if count >= limit then
    -- waitTime = today
    waitTime = 86400 - (now - today)
elseif now - lastTime < interval then
    waitTime = interval - (now - lastTime)
end

-- 更新发送次数和最后发送时间
if waitTime == 0 then
    count = count + 1
    lastTime = now
    redis.call("HMSET", KEYS[1], "count", count, "last_time", lastTime)
end

-- 返回结果
local result = {waitTime, cjson.encode(redis.call("HGETALL", KEYS[1]))}
return result
