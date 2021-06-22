### 실행 명령어
```$docker-compose up```

- redis 사용 안하면 go run main.go 로 실행 가능 (router.go에서 tokendb redis 주석처리하고 temptokendb 주석 해제)

### jwt login 구현
기본 컨셉
- 로그인하면 15분 동안 유효한 accesstoken과 accesstoken 만료되었을 떄 재발급 시에 사용할 refreshtoken 생성

#### auth.go 
login 함수 
- db에서 id, password 확인
- db에 값 있으면 토큰 생성
- 토큰 redis에 저장
- 생성된 토큰 반환

logout 함수
- request header에서 토큰 추출
- 토큰 db에서 삭제

tokenrefresh 함수
- access token은 만료되고 refresh token은 살아있을 때 access token 재발행하는 함수
- request body에서 refresh token 추출 (원래 body에서 받는 건지 확인 필요)
- refresh token 유효한지 확인하고 유효하면 기존 토큰 삭제하고 새로운 토큰 발행 및 반환

#### tokenhandler.go 
createtoken 함수
- access, refresh key 갖고 jwt 토큰 생성

extracttokenstring
- request header에서 access token string 추출

gettokenfromtokenstring
- 토큰 변조 여부 확인(사인 방식 같은지 확인)
- token string을 token으로 변환

checktokenvalidation
- 토큰 유효성 확인

extractuseridanduuid
- 토큰에서 userid와 토큰 고유 uuid 추출(tokendb에서 확인할 때 사용)

#### tokenmiddleware.go 
tokenauthmiddleware
- router에서 함수 실행하기 전에 실행되는 함수 
- 라우터 함수 실행 전에 유효한 access token 갖고있는지 확인(인증 필요한 곳에 접근할 때 사용)

### 회원가입 구현하기

token.claims / token.mapclaims 둘 다 valid 체크 해야하는지?? 

access token, refresh token 다 request에 들어있는지?? access token만? 
