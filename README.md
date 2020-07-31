# Twelve Factor App ENV buddy

App ENV vars
* [AppName]_[AppEnv]_VAR_NAME

Directory Structure
```
.tfa
|   tfa.config
|   tfa.sqlite
|   tfa.sh
```

###Commands
* tfa create [app name] [env]
* tfa delete [app name] [env]

* tfa set [app name] [env] [var_name] [value]
* tfa remove [app_name] [env] [var_name]

* tfa commit

* tfa view

#### Future Commands

* tfa disable (app, env and var level) 
* tfa enable (app, env and var level)
