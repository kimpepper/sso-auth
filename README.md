# SSO Auth

A proof of concept for AWS SSO CLI login.

## Usage 

### Configure AWS CLI for SSO.

#### Automatic

```
aws configure sso --profile=my-profile
SSO start URL [None]: [None]: https://my-sso-portal.awsapps.com/start
SSO region [None]:ap-southeast-2
```
...and follow the remaining prompts.

#### Manual

Add profile configuration in `$HOME/.aws/config` for each account you want SSO access:

```
[profile example]
sso_start_url = https://<start url prefix>.awsapps.com/start
sso_region = ap-southeast-2
sso_account_id = <account to access>
sso_role_name = <role to assume>
region = ap-southeast-2
output = json
```

### Test the credentials

```
./sso-auth --profile my-profile
```
