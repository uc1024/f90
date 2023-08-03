-- Allow 当前动作是否是被允许的

-- KEYS[1] key 主key
-- ARGV[1] uinx_now 当前时间
-- ARGV[2] limit_tiem 允许的访问间隔
-- redis-cli --eval ./test.lua 1
--
--   @Last visit            上次访问时间
--   @Overclock             被限制时调用的次数
--
-- @ return
-- 1 允许访问
-- 2 间隔限制

local unix_now = tonumber(ARGV[1])

local last_visit       = "last_visit"
local overclock        = "overclock"

local if_have = redis.call("HEXISTS", KEYS[1], last_visit)
if if_have == 0
then
    -- init
    redis.call("HSET", KEYS[1], last_visit, unix_now)
    redis.call("HSETNX", KEYS[1], overclock, 0)
end

-- 访问间隔限制
local time_last_visit = tonumber(redis.call("HGET", KEYS[1], last_visit))
if time_last_visit + tonumber(ARGV[2]) > unix_now
then
    -- 访问间隔限制中
    redis.call("HINCRBYFLOAT", KEYS[1], overclock, 1)
    return 2
end

redis.call("HSET", KEYS[1], last_visit, ARGV[1])
return 1
