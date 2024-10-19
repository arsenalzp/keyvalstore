package main

import (
	"crypto/tls"
	"net"
	"os"
	"testing"
)

const CRL_FILE = "test.crl"
const CRL_DATA = `-----BEGIN X509 CRL-----
MIIBvTCBpjANBgkqhkiG9w0BAQQFADB3MQswCQYDVQQGEwJVQTEVMBMGA1UECAwM
WmFwb3Jpemh6aGlhMRUwEwYDVQQHDAxaYXBvcml6aHpoaWExETAPBgNVBAoMCEdP
S0VZVkFMMREwDwYDVQQLDAhHT0tFWVZBTDEUMBIGA1UEAwwLZXhhbXBsZS5jb20X
DTIzMDYwMzEzMTQxNVoXDTMzMDUzMTEzMTQxNVowDQYJKoZIhvcNAQEEBQADggEB
AJbUat/LI8hnLXM7YJQH4ceQWPylkueOfXQYUpeNUuH52p9SNFZaqBrsvmr+sv2T
HolQ/Fy3Yhd0VavRDGrIav3h9RefJ7pY4nCiuv0qJExcLEx1pehYLI+XnnVvNXPA
Yb/RoBjWu/b+XY/VY+O19iqzpYyyP//Mg6q2GIWzORBoIt5sZDUm70g2fWB0oZxy
lAnoLlv2b+1/wiiu+AzYu9DPgsnOZv5f/QqJdvdEk95wE0e3l+BmsWFvrLhLq1oY
f7FB4HoJZJDEM+UmMLSzLJVrBYn0Y8qBsXyl0wyUxiX6pRDIG9+zyLmS9j6ZpPvU
gxSPtH55WVnnZ5k5SEf/Vvg=
-----END X509 CRL-----`

const ROOTCA_CERT = `-----BEGIN CERTIFICATE-----
MIIDzzCCAregAwIBAgIUXS7P6grVGLNlKO6lXUzmzsfCOAswDQYJKoZIhvcNAQEL
BQAwdzELMAkGA1UEBhMCVUExFTATBgNVBAgMDFphcG9yaXpoemhpYTEVMBMGA1UE
BwwMWmFwb3Jpemh6aGlhMREwDwYDVQQKDAhHT0tFWVZBTDERMA8GA1UECwwIR09L
RVlWQUwxFDASBgNVBAMMC2V4YW1wbGUuY29tMB4XDTIzMDYwMzEzMTQxNVoXDTMz
MDUzMTEzMTQxNVowdzELMAkGA1UEBhMCVUExFTATBgNVBAgMDFphcG9yaXpoemhp
YTEVMBMGA1UEBwwMWmFwb3Jpemh6aGlhMREwDwYDVQQKDAhHT0tFWVZBTDERMA8G
A1UECwwIR09LRVlWQUwxFDASBgNVBAMMC2V4YW1wbGUuY29tMIIBIjANBgkqhkiG
9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1ns3CjKK2uaGyocsb8z28WMCuYzc5UDPiCG7
zhP5KfILJDadZtev/T3TGCj6SR2DHk8SAvwATNnZSfytAc7zCPncFxT36iTmU0jZ
Jp2TaVcdSoAtQ+YtyOVNdJXjk69yRQjLOnrEGvSe8kcMaeK6YXPhJtlxv4tfdgO8
5Ky41gb+RvDS/CoCwQiEwWdpQ36qDr+69Jpm67t38sxZiKlabhZ9xdeFpLUfB6GW
9Fk+99sIOn9WNUG439J6x9uZDcKPcw/ylhxnDQiOHhH4ogqW5NQn8oB7Wgk5bMhj
eA2VhazxGE4+S2h3B43Z20jeZoEadN3USljxQCKJWyDErivTuQIDAQABo1MwUTAd
BgNVHQ4EFgQUhmOhP13toMflJuf6pasKpB3crSkwHwYDVR0jBBgwFoAUhmOhP13t
oMflJuf6pasKpB3crSkwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOC
AQEAJ1uxshIrLdjXb1L8S0pe6I2JddI3dRiN5EQOhbUS//kRr/GXtbMzeukHSZhV
pgKGshcF6DFlmulq0ciMLWlzTsifEx5vDCXFYvZOR3fseRD5Ul8Mr9+UeEeJLrPD
DJj4zkQh4lRCWwvRT3xWTAqyOA3PKSPcFu1yaVZE4cvkPsPJJ+QySJBi1P+lIIhq
FQVYXW4V2rJoI23Ts5g+z4Xv062GQ/iN5Lqmekv8lwhkIBBQC4AW1jJZ2Q4fMeGw
e8utdDncC3NTtJRrqf2ozZKWo9OfQZL0qe2gxM5c1ilRirQtJa7y6aXDrAvMRrco
+ipQSZIG3GkWYl3xxZHBGeP9yA==
-----END CERTIFICATE-----`

const ROOTCA_KEY = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDWezcKMora5obK
hyxvzPbxYwK5jNzlQM+IIbvOE/kp8gskNp1m16/9PdMYKPpJHYMeTxIC/ABM2dlJ
/K0BzvMI+dwXFPfqJOZTSNkmnZNpVx1KgC1D5i3I5U10leOTr3JFCMs6esQa9J7y
Rwxp4rphc+Em2XG/i192A7zkrLjWBv5G8NL8KgLBCITBZ2lDfqoOv7r0mmbru3fy
zFmIqVpuFn3F14WktR8HoZb0WT732wg6f1Y1Qbjf0nrH25kNwo9zD/KWHGcNCI4e
EfiiCpbk1CfygHtaCTlsyGN4DZWFrPEYTj5LaHcHjdnbSN5mgRp03dRKWPFAIolb
IMSuK9O5AgMBAAECggEACFq0zmIhG/0yP2XMAo/Uh5sAcmtM+GYRn18XAmQvRqW9
rtusvmKE4aOcmxqPgzUR1Cam0DRyK4wpxVLBB2rA9zbKTqDbtkRZL2HNSY6GPSpu
QuZivPwUiPh7j6C+BFTnac55og3BoT4piczgLCivc1fQ65M5NcmPPpjPLQYYTgfV
4khLvCv114nH5xtavueJ0RENUj++ijl0huDj9u8zWdnhAiNQswi7hmy/H6QKPwxH
IEjN284yU/Cn2ek1MXGQiz6L8VrYDAztKB9rfks8n2gF0jCaD8KzvQ5fduzGNyrp
WL4aGOJMSVXQmpK9KjctyBE7ftYijb53wM6z0boF4QKBgQD/B9iZIJcmWHrBPxg+
/sra5Iv3w68zR2Ycdd9XQpvo1tG4d0emvHjjqsvveva4ja3RPltIRLl+LQcjICo5
4484D+owq7TXCEQZRask/FSw/CFhxDsnFoUuUkWcqocSMtBs1RNX6YiEX6Z45IUU
/oa+WTuprae/LHjxXiagAJj2EQKBgQDXS+m3sTv2qGwf1Ore7dpfD9c2KV2lniqc
drEBMCYilkSPJt4CHzkKu/OoygWZIkuvBMWyq45okIyJEDU8GbP7Xu0ttKx3r22q
N3AqpbZoagzU2d1PNpPp1WEA/OK8MMb0Q6+/MqBgM+pLCbv3PZLRHVyp0MBmT9Sq
fPGoFQ67KQKBgQCoa3YPKganvCbVF+XbNEii5evJY+F/69bzVKR56/MqgTNerucS
pcTwFQs4y+vDVU1Esfl2cGxPd00PVV6NfEpIq7ntCngSydKvHeM4Oat0dg1Vk1G9
LpIlVQ/Dtpoc2pHqTYzIseEGCmTV0ZRRmQVDD4rnM1dkWOpF+/dmEv7xAQKBgDyx
Rsqk6P2IzypOEIQV00iXI5k+IgstQl/nSdDG7Qk1CVC9qTo56Q+wmmjLwrY/p8xH
8R+EI1ow3Z1J92fg6w8C3KPU2gXHa4ffpvwuyPQ4aTOb0zqgbSQvJfBsWdKpgXyc
lC+3KuTT3cmXjeiN8BSJTXUFxydQe+gv3sP+Y6+5AoGBANYRC7wNQ82FI1zSGFDU
SGKCFz5jNgzOEAZf++lBGGWfiDzOINkoZtk+Ng3C7V8mMJbiEcRUWVPoXRXvOzWk
eTTIZv4OWchr/irigkurzT/QX/wjV0Z6tSrmGzTe6e0KmxK5+kI2NV/0SVLaTOik
7P63RT8/zsCyYneALEjFt5D5
-----END PRIVATE KEY-----`

const ROOTCA_SRL = `42B12FCA51930F6B395687322A72563337068ED5`

const SERVER_CERT = `-----BEGIN CERTIFICATE-----
MIID+TCCAuGgAwIBAgIUQrEvylGTD2s5VocyKnJWMzcGjtQwDQYJKoZIhvcNAQEL
BQAwdzELMAkGA1UEBhMCVUExFTATBgNVBAgMDFphcG9yaXpoemhpYTEVMBMGA1UE
BwwMWmFwb3Jpemh6aGlhMREwDwYDVQQKDAhHT0tFWVZBTDERMA8GA1UECwwIR09L
RVlWQUwxFDASBgNVBAMMC2V4YW1wbGUuY29tMB4XDTIzMDYwMzEzMTQxNVoXDTMz
MDUzMTEzMTQxNVoweTELMAkGA1UEBhMCVUExFTATBgNVBAgMDFphcG9yaXpoemhp
YTEVMBMGA1UEBwwMWmFwb3Jpemh6aGlhMREwDwYDVQQKDAhHT0tFWVZBTDERMA8G
A1UECwwIR09LRVlWQUwxFjAUBgNVBAMMDSouZXhhbXBsZS5jb20wggEiMA0GCSqG
SIb3DQEBAQUAA4IBDwAwggEKAoIBAQDbZD2YCSN8leJA26p9R+CFX2wt6KjGX8Gn
kzQrRXl2mbgNjHmQnRjdEDskzdSFlDp/ZIdVB2sOIYqE8V/xY3EQt6tkNUiew218
QjRAbFdN2YcucywCdiNH6fZ8j6gQKEU8zb/DyZxk+BNSYvV3qEBogHXRsnOQD2V0
2vvagfhTO5kaDtP4Upq3Ma0DB4dMZPZzuxBLhUUUfZvdgq3ZCfc+LiA1lJNy+inv
0A+atR/wAnWQpO7yDyO/Z1M8MDpDX5Z+Z8Y14lwGO1EjB1BY0soQUL75kbPw0U06
jKs9epH8SzvdBlVBjjdWIa/hA71Hz19jWtqon5idBjxg1TYBfF+pAgMBAAGjezB5
MDcGA1UdEQQwMC6CEnNlcnZlci5leGFtcGxlLmNvbYINKi5leGFtcGxlLmNvbYIJ
bG9jYWxob3N0MB0GA1UdDgQWBBT94qFds3SNSDVDTS1Tnf9V7cEwXzAfBgNVHSME
GDAWgBSGY6E/Xe2gx+Um5/qlqwqkHdytKTANBgkqhkiG9w0BAQsFAAOCAQEA1EPo
H/4f4RlSyRAZnJjB9fwOGekwIOuCq2beIf26Gz8OnJSA6DBQqbMq3uDDSGs564nG
vGQ/c9H1pbwamo5Y+TzFetalgOdg4VhkTLfx+vTnM//N6BDRiXZyhF3AUD+wzQqp
eomwcG3lrpJlhWAln/CtH4x/6R3ADZERfavqbm0UErkTfihyLnIOYaQg48lb/j0M
jQ3i1jUBXxl6V3Z+JUkjgfYP/WoQzJaXQYVYXzLyQCGFLLIM+KKX3QM7ti5w/oF1
v2Om0EAgtlJ/5+Nhw39Cl/IPuMbm0mQhc6w0xzi96cgU1owszOLazZdxawgtP47H
nNJNZljxPvAk4unhuw==
-----END CERTIFICATE-----`

const SERVER_KEY = `-----BEGIN PRIVATE KEY-----
MIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQDbZD2YCSN8leJA
26p9R+CFX2wt6KjGX8GnkzQrRXl2mbgNjHmQnRjdEDskzdSFlDp/ZIdVB2sOIYqE
8V/xY3EQt6tkNUiew218QjRAbFdN2YcucywCdiNH6fZ8j6gQKEU8zb/DyZxk+BNS
YvV3qEBogHXRsnOQD2V02vvagfhTO5kaDtP4Upq3Ma0DB4dMZPZzuxBLhUUUfZvd
gq3ZCfc+LiA1lJNy+inv0A+atR/wAnWQpO7yDyO/Z1M8MDpDX5Z+Z8Y14lwGO1Ej
B1BY0soQUL75kbPw0U06jKs9epH8SzvdBlVBjjdWIa/hA71Hz19jWtqon5idBjxg
1TYBfF+pAgMBAAECggEAMc6nlyEX/9xNQdk88vVC+TSJLEECNAsgPWmTcfFzWmQa
n2TRYm3v65wYOUpLYcodn7dUbA7jlJEzz+u2kug3DosMK5NXOcf3Trr+/tM53NAy
Ou7lwmdMqjJpBa1dg9GIqn2xeAMI8PlK9azGupQljzP+y40eZEnCiE2A2QIhvM9C
qfC1C6dPzuHXN0SHvry4z+3Frw0C00swDzHALT+w47cfinm7Fxyn3B58qVqdQjjC
v/f1B2ku7chRmKQAVGBMOaCtVUa1Dj7s7Q0yuq+A36ob/OGLTGfxOf/rK+fSh3Ov
BNyTAyKTGLJRnDjbWc2eDpckX3zbwYsh3oNM+8lwWwKBgQDeuFApma5DjUrkmPxW
QIqogoDOnMHpnHVrbosed235ygvqCML/A/5iSSSn7pg/isteyIeKBOiTEK5tyq2O
8hgJntm6IJ7zNPhPSPLukT/1lQrFp6QI0ZHXWJjKltO4Y70v0sYngh2lOBm8BoBY
0kgJ3CIwxEXyk5pSzhwb0Ty0MwKBgQD8LJstpUILBgzXR0QeFtkU/Nt+DZdAcz1o
9RL66RDtsczDFQs6PgCW5tD4MnhQ6sSiitljLqZ+qxrPLfMCB4hAeJzjWTwtXN8x
KzrQvklx+5IH0yvFQkZ2B6YI1IGFygkoHphhBtxJARMjyfP4Yp6oxqpIyZlz9dfN
6XSsiGwgswKBgHqMlmddjZrT+xqv52EaYHF6ZZ/Kd5Swp5d2mwwnkRb6CvY63fju
XKH+NzJEQffsyhycYKAcNVD+w8vb0wYtxfY9NvaIjo2qXttZe3qz56qc2PGLXeIQ
VpxUvrXyqgryrp3K74e7u842gUqJlUPKaSMrwpBs30Qr3aWkjajsx+crAoGAJ9kU
nF3k1cEa/lmwleCeZQaf2IdlQzXymkc/vI5fsm/KH3mP0KBDj5ThqJaxFHhEojq2
p0mT3ahEEED+iW+PREDK6dIMBE8MpcRjAuFO0cgjB0GDRSR35ebHgdWyseV/FOvg
wFRJMvAMijc7aiCLWbgq6F2S9hP/Cfa+DRVxoKkCgYBDf0stWctIt6LRFLgVaF3k
ewQP6696GGyFY4qFyLQOE/jKXG9P1bIGh4CLqGseSvZ3NLmosfO6HhpnaQ6HdaLX
4PnebRd9ptZU9h05JeHzhgdxfpDTREttv9qLqkEfTav7UmTmsyW72kf59U7WuquZ
jWAtsbuEr68235aZA+Krug==
-----END PRIVATE KEY-----`

const CLIENT_CERT = `-----BEGIN CERTIFICATE-----
MIID+TCCAuGgAwIBAgIUQrEvylGTD2s5VocyKnJWMzcGjtUwDQYJKoZIhvcNAQEL
BQAwdzELMAkGA1UEBhMCVUExFTATBgNVBAgMDFphcG9yaXpoemhpYTEVMBMGA1UE
BwwMWmFwb3Jpemh6aGlhMREwDwYDVQQKDAhHT0tFWVZBTDERMA8GA1UECwwIR09L
RVlWQUwxFDASBgNVBAMMC2V4YW1wbGUuY29tMB4XDTIzMDYwMzEzMTQxNloXDTMz
MDUzMTEzMTQxNloweTELMAkGA1UEBhMCVUExFTATBgNVBAgMDFphcG9yaXpoemhp
YTEVMBMGA1UEBwwMWmFwb3Jpemh6aGlhMREwDwYDVQQKDAhHT0tFWVZBTDERMA8G
A1UECwwIR09LRVlWQUwxFjAUBgNVBAMMDSouZXhhbXBsZS5jb20wggEiMA0GCSqG
SIb3DQEBAQUAA4IBDwAwggEKAoIBAQDbsh36oFWJpdSrCqXjST6VpWxZ2JyBoOTC
Zs2N+W5jN99iDgCyArhKTubI5Js1z40jnXEY59bUGomQ9BfKcLZrksIFjXiPbNo6
eJdAvwR9TjRr5Uj2sBMy4okmLwt2Z7j17FOV647vCfqr0scEqYbGwHJqA4Jod3pW
W6v2WpBSUiQEXvCzVi42847F2XRyoYiz+kdXDZziP6CFfOWyqoD0pl/jePpvHIoN
jeKWMIlnQQPwELxtVXuCmX2LR7p1KUMBgfEEZtd4hJNY/9/NB6xAa9nohdm6FcP7
0eRJ1h2IQTZguzBp91Ff/bBlHTVx26HqqO6j8aQwMOYCn8JHTq/JAgMBAAGjezB5
MDcGA1UdEQQwMC6CEnNlcnZlci5leGFtcGxlLmNvbYINKi5leGFtcGxlLmNvbYIJ
bG9jYWxob3N0MB0GA1UdDgQWBBTHmdnYSnqLxBUw1FkJukrxgFPn8TAfBgNVHSME
GDAWgBSGY6E/Xe2gx+Um5/qlqwqkHdytKTANBgkqhkiG9w0BAQsFAAOCAQEAvRqE
YWF08t9BQubIwVGoERV+CSm84xsTdZgHEf14glRNd1+k1gaNZS6K7+Da81t1ATFc
s4z4nJLGZVUwgtMraMpiP6X6j5/jWKX5Zfzd2iCcFfH6uAmTP8ABXvAOlaStGE1V
jgRckCuxhfTUt6gJAfBfo0jt/WmIXJiosNGvmT6xtBZPTtXqZIa2Qc9vdYbtWbkn
jdBvjKoh1dyhVPyn19I3y3Z3CXffLx/0lw2gcr3mvX5LxqoyxcukPCSSoKyQOKC5
G8yBJOq42jJ1oDwjTEa/7CeZ7d6dOhXUkJhkWX5O/6WsXFL0M1IbUnRu3fe3IOJs
pf2EWBfCDqJvaR75zQ==
-----END CERTIFICATE-----`

const CLIENT_KEY = `-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDbsh36oFWJpdSr
CqXjST6VpWxZ2JyBoOTCZs2N+W5jN99iDgCyArhKTubI5Js1z40jnXEY59bUGomQ
9BfKcLZrksIFjXiPbNo6eJdAvwR9TjRr5Uj2sBMy4okmLwt2Z7j17FOV647vCfqr
0scEqYbGwHJqA4Jod3pWW6v2WpBSUiQEXvCzVi42847F2XRyoYiz+kdXDZziP6CF
fOWyqoD0pl/jePpvHIoNjeKWMIlnQQPwELxtVXuCmX2LR7p1KUMBgfEEZtd4hJNY
/9/NB6xAa9nohdm6FcP70eRJ1h2IQTZguzBp91Ff/bBlHTVx26HqqO6j8aQwMOYC
n8JHTq/JAgMBAAECggEAFierMEgU+DGZ9bm4KuiLACpTd+gJOGVSTGxzlDqwMB9F
Tq0c0tbFYPD+AwpSwKVylUHeUuWmW3NlphGHiKm/K6/8EvVGUChpBXZ9wlDBEiXd
0Xeo2P++n+YKcKhT3pftJhe0Ai1kF6UI+2ThHw49gjOMFjbOYtyRoL3T5J8TaMmh
55lYo9pH9uZN+cKa5U63v+j4lEleYXayrNNOxU/Gme6tCE+BcMdr/qilXz0+o+dO
U9aNZVI2nkxNCGUOnd99+FQ8ZLKh8XGWIZZ6RSorcghNd86U6RAtHc4d1OuQ3Sc9
tKa2hRxaMxGAY3Srp3nMbsxCOhCuTCLQYJuQfAqg1QKBgQD4n9hiybiSjQIqd+q1
QiHL9O4pkgJDTSOhU5o7L/wqTAO23iOTp4TtMYSijyb6thpWGH/N2//T40wsCm0k
93CcWYLkg/+EYVVoxh4GMiS4aAvmR8lfDjuFqJGDfhqQPPKyzzjvtAYseoUp7OT4
Figp54trWzpfFQ8xqsjDFwePHwKBgQDiNpN7i2BWFrNX1EgzDKUxHu1rTqwJLMr1
LlKMbZaFExt562YlCQVlivBXnyBCw/p5Mg5O6B3DLowiQiNnyVPyKGUgYA8PsEiN
4joOJMQsWhHNvaFz++AzopJk8GY6GiXY+74rVfYhAa53K8EgB7S2PH5tTwiFpUuY
jcLjx4isFwKBgQDB39vcBRNR7HVw6nvzBnPWWNPTRNFQ6/lJ1yig2OVZklcfJZA2
lt4YHJIiNWEfBhv5YTdgLxsKfueqPCMqPW3p7f8c9TWuZDw27K8DA90Qk8obs4T7
A900d+Oo1xAdw/k5qE/s08QwsQXgUKOoNZbyPmXAvK4C8SgdAeF2CCJT3wKBgQDf
8SV55dW+BAURis7a8sbKZRKm66A2CQj3Rh9kc8zR+sN1pAtf2JlmF/Cs3ZQDZJ4e
wuYVSYbFRdxmwdDpGw8mqMTMEyx13I9HHtFYVR97xMLhSbx+5LfkhimlEbQyCtaz
Ay0VG6lorZB423D584b77dE/B0GphKTc5mIsOslbiwKBgQCK3sdb7ozGccYGu/WA
XzF+uotcZvrUB9nSU4pX5st6DbsOEXVv8spFv4dMONx8u0+20n0Al3UIfw706YDr
sWvnfseTpElMgncH5/ek+qfjOUK5NxsFcn4Mxf0/9MZLnDlSvSGwZkVENqQ/jPuR
cZmdS6IeVueSRnc9yOZET+b8oA==
-----END PRIVATE KEY-----`

func TestStartServer(t *testing.T) {
	err := prepareCRL(CRL_FILE, []byte(CRL_DATA))
	if err != nil {
		t.Errorf("unable to create CRL file, %s\n", err)
		return
	}

	defer cleanUpCRL(CRL_FILE)

	srv := Server{
		CrlPath:        CRL_FILE,
		ServerCrtData:  []byte(SERVER_CERT),
		ServerKeyData:  []byte(SERVER_KEY),
		RootCACertData: []byte(ROOTCA_CERT),
		Port:           9999,
	}

	lsnr, err := srv.Start()
	if err != nil {
		t.Errorf("unable starting server, %s", err)
		return
	}

	defer lsnr.Close()

	if lsnr == nil {
		t.Error("unable starting server, listener is nil")
		return
	}
}

func TestSSLConnection(t *testing.T) {
	err := prepareCRL(CRL_FILE, []byte(CRL_DATA))
	if err != nil {
		t.Errorf("unable to create CRL file, %s\n", err)
		return
	}

	defer cleanUpCRL(CRL_FILE)

	srv := Server{
		CrlPath:        CRL_FILE,
		ServerCrtData:  []byte(SERVER_CERT),
		ServerKeyData:  []byte(SERVER_KEY),
		RootCACertData: []byte(ROOTCA_CERT),
		Port:           9999,
	}

	lsnr, err := srv.Start()
	if err != nil {
		t.Errorf("unable starting server, %s", err)
		return
	}

	go func(l *net.Listener) {
		conn, err := (*l).Accept()
		if err != nil {
			t.Errorf("error establishing connection: %s\n", err)
			return
		}

		tlsConn := tls.Server(conn, srv.GetTlsConf())

		err = tlsConn.Handshake()
		if err != nil {
			t.Errorf("error establishing secure connection: %s\n", err)
			return
		}
	}(&lsnr)
	defer lsnr.Close()

	clientPEM, err := tls.X509KeyPair([]byte(CLIENT_CERT), []byte(CLIENT_KEY))
	if err != nil {
		t.Errorf("unable to parse clients certificate or key: %s\n", err)
		return
	}

	clientConn, err := net.Dial("tcp", "localhost:9999")
	if err != nil {
		t.Errorf("unable to establish secure connection with the server: %s\n", err)
		return
	}

	// create secure connection for client
	tlsClientConn := tls.Client(clientConn, &tls.Config{
		ClientAuth:         tls.RequireAndVerifyClientCert,
		MinVersion:         tls.VersionTLS13,
		Certificates:       []tls.Certificate{clientPEM},
		InsecureSkipVerify: true,
	})

	err = tlsClientConn.Handshake()
	if err != nil {
		t.Errorf("unable to establish secure connection on the client side: %s\n", err)
		return
	}
}

func cleanUpCRL(file string) error {
	err := os.Remove(file)
	if err != nil {
		return err
	}

	return nil
}

func prepareCRL(file string, data []byte) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	return nil
}
