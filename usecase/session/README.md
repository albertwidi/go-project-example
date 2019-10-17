# Session

Allowed multiple session running on multiple device

For example user `A` might have `android` and `ios` device, the user might use the same account to login into two different device and the session will be handled separately.

## Session Storage

Session storage is using `redis`

## Multiple Session

To ahchieve multiple session, `HASH` is used in redis, and the session is formed in this fashion:

`usersession:{userid} {hash_sessionid1} {hash_sessionid2}`

So it is possible to track how many session a user had and allowed us to delete all the session and force user to logout from all device.

## Tracking Expiry

To track the expiry and allow each single `HASHFIELD` to have its own `expiry_time`, `expired_at` field must be added to the session.

When a key is `retrieved`, program will check whether the `session` is expired or not. If the `session` expired, then the program will force the `user` to re-login and dispatch a job to delete the `session` from `HASHKEY`

## Downside

### State Sync

The state between multiple session instances needs to be synced, for example if a user status is changed.

### Session Expiry

User will only have 1 expire time. Once 1 key is expired, the other key will get expired too. The management of the `HASHFIELD` is become more complex than simple `session`.

Many keys might be dangling when user is no longer active. Because we only have one time to expire, the expire time will be far longer than the usual `session` key, or maybe expire will never happens. All key `deletion` and `renewal` is based on user action.

To delete the remaining `dangling` keys we might need some activity tracker, and by using cron to delete the `dangling` user session one by one.
