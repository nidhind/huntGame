# huntGame
Server for huntGame

### Environment variables to export


| Name | Description | Default Value |
| ------ | ------ | ------ |
| `SENDGRID_API_KEY` |Required for sending emails |`""` |
|`RESET_PASSWORD_URL` |Redirect URL used in reset password email link |`""` |
|`RESET_PASSWORD_UPDATE_REDIRECT`|URL to which the server will redirect after the validation of `reset_token` in the email link, preferably the UI route|`""`|
|`JWT_SECRET`|Secret key used to sign `JWT` tokens|`""`|
|`GO_ENV`|Running environment for APP. Possibly values are `development`, `production`|`development`|
|`PORT`|`port` on which the server will listen and serve|`8080`|
|`RESET_PSWD_TEMPLATE`| Absoulte path to Reset Password Email Template file. If not set then default in-build template will be used|`""`|
|`AUDIT_LOG_BASE_PATH`|Absolute path to the folder in which the `audit log` file must be created.WARNING: If not set then tmp folder will be used and the data will be lost when the instance restarts on normal systems|`/tmp`|
|`AUDIT_LOG_FIlE_NAME`|Filename to be used for the `audit log` file|`audit-answer.log`|
|`GIN_LOG_BASE_PATH`|Absolute path to the folder in which the `gin log` file must be created.WARNING: If not set then tmp folder will be used and the data will be lost when the instance restarts on normal systems|`/tmp`|
|`GIN_LOG_FIlE_NAME`|Filename to be used for the `gin log` file|`gin.log`|
|`LEADER_BOARD_LIMIT`|Number of entries to return for leaderboard api |`20`|
