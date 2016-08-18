* Detect similar behaviour from different IP addresses
* Check distance (in feature space) between sessions from the same subnet
* Close session: take 1m sessions, calc max idle period in each session
* Milliseconds in date
* Add source port to session key to distinguish users behind NAT



# Features
* Every API is a set of features:
  * delay of the first request since the beginning of the session
  * average delay between separate API calls (request time minus prev. request's response time)
  * chaining (have referrer been requested?) 1,0
  * number of calls per minute
  * track order of calls (how?)
  * track delay between calls to different APIs (temporal patterns)
* Track set of cookies for every API call
* Requested static files
*

* Session general features
  * session total duration
  * session start GMT hour
  * session start in client timezone
  * start minute (to detect attacks started in round hour time etc. 4pm)
  * Mobile device?


unusual browsing time (for client's time zone)
each API (by content-type):
number of calls per session
average delay after the last request (any)
average delay after last same request
track content-type (to distinguish images, apis, pages)
session started near round hour time (like 3am, 12pm)


* Ideas

  * do not append features from separate request into a single feature-vector
