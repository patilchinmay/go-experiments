# HTTP Client with Auto-refreshing tokens

Demonstrate creating an HTTP client that is capable of refreshing its tokens.

We will use the OIDC-Oauth2 client_crendentials flow for token purpose.

- [HTTP Client with Auto-refreshing tokens](#http-client-with-auto-refreshing-tokens)
  - [Keycloak](#keycloak)
    - [Start](#start)
    - [Configure](#configure)
    - [Test](#test)
  - [Run the program](#run-the-program)
  - [Explanation](#explanation)
  - [Reference](#reference)

## Keycloak

We use Keycloak as the auth provider.

### Start

Start the keycloak as a container.

```bash
docker compose up -d
```

### Configure

Configure the keycloak realm and client:

1. Visit http://localhost:8080
2. Login using admin credentials.
3. Create reals `test`.
4. Create client `test-client`.
   1. Set `Client Authentication` to On
   2. Set `Service Account Roles` to On
   3. Set appropriate home, root, valid redirect and web origin URLs (for testing we will use `*` wherever applicable).
   4. Copy client id and secret from the `Credentials` tab after the client is created.
5. Reduce the token lifespan for testing purpose at http://localhost:8080/admin/master/console/#/test/realm-settings/tokens
   1. Set access token lifespan to 1 minute.
   2. Having smaller than 1m can lead to issues.
   3. The client will refresh the token when it expires (or is about to expire), but due to the short expiry time and your request interval, client may renew the token (even though it is not currently expired).

### Test

Test client_credentials flow by logging in for service account:

```bash
curl --location 'http://localhost:8080/realms/test/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'grant_type=client_credentials' \
--data-urlencode 'client_id=test-client' \
--data-urlencode 'client_secret=<test-client-secret>' \
--data-urlencode 'scope=openid email profile'
```

It's important to note that the client_credentials flow typically doesn't return a refresh token. This flow is designed for machine-to-machine communication where the client itself is the resource owner. In this scenario, when the access token expires, the client is expected to request a new access token using its credentials rather than using a refresh token.

## Run the program

```bash
go run .

2024/07/28 17:55:11 req no: 0 response: 200 tokenID: 4d046800-075c-4d22-839c-c5aa71f68293
2024/07/28 17:55:41 req no: 1 response: 200 tokenID: 4d046800-075c-4d22-839c-c5aa71f68293
2024/07/28 17:56:12 req no: 2 response: 200 tokenID: 58278c11-24ef-4fd9-a403-46367378eb28
2024/07/28 17:56:42 req no: 3 response: 200 tokenID: 58278c11-24ef-4fd9-a403-46367378eb28
2024/07/28 17:57:12 req no: 4 response: 200 tokenID: f0df2a07-98fd-4792-b449-5d2e18e2cf58
2024/07/28 17:57:42 req no: 5 response: 200 tokenID: f0df2a07-98fd-4792-b449-5d2e18e2cf58
2024/07/28 17:58:12 req no: 6 response: 200 tokenID: 978002e1-86c7-4822-8ed1-575c3c9eb797
2024/07/28 17:58:42 req no: 7 response: 200 tokenID: 978002e1-86c7-4822-8ed1-575c3c9eb797
2024/07/28 17:59:12 req no: 8 response: 200 tokenID: 2c25ee1c-053a-492d-86b4-e3908e5d2520
2024/07/28 17:59:42 req no: 9 response: 200 tokenID: 2c25ee1c-053a-492d-86b4-e3908e5d2520
```

## Explanation

The `access_token` lifespan is 1 minute. Our program sends a request every 30 seconds.

So, ideally, our client should use the token for at least 2 request (30s * 2 = 1m). It should get a new token after every 2 requests.

This can be verified from the logs above from the `tokenID`, which is same for at least 2 consecutive requests.

## Reference

- https://darutk.medium.com/diagrams-and-movies-of-all-the-oauth-2-0-flows-194f3c3ade85
- https://www.youtube.com/watch?v=MeXzVS4QZ4Q