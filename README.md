# srcdsLogParser

this is log parser at JavaScript for valve`s officials server (srcds) : CS:GO, HL2, DOD:S, CS:S and etc

### who to use

1. download **parser_v2.7.js**
2. add to your project like
  
  > const parser = require('./parser_v2.7.js')
  
and that is it

it will return promise object, use it like that

>    parser( line-to-parse )
>        .then( msg => { ... })
  
 object will contains 
 event: ...
 eventTime: ...
 player {
  name: ...
  id: ...
  steamID: ...
 }
 
 and etc
 
 ### TO-DO
 more desctiption and get it to NPM :)
