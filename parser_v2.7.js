let fs = require('fs');

function teamNameFn(line){
  
  return new Promise(function(resolve,reject){
    let result = line.match(/^Team playing ["](.+)["]: (.+)$/i);

      if (result) {
      let outResult={};        

        outResult.event ='teamName';
        outResult.team = result[1];
        outResult.teamName = result[2];

        return resolve(outResult);
       } else {
       return resolve(undefined);
       };
  });
};

function attakedFn(line){
  
  return new Promise(function(resolve,reject){
    let result = line.match(/^["](.+)["] \[(.+) (.+) (.+)\] attacked ["](.+)["] \[(.+) (.+) (.+)\] with ["](.+)["] [(]damage ["](.+)["][)] [(]damage_armor ["](.+)["][)] [(]health ["](.+)["][)] [(]armor ["](.+)["][)] [(]hitgroup ["](.+)["][)]$/i);

      if (result) {
      let outResult={};        
        outResult.event='attacked';
        outResult.playerA = parsePlayer(result[1]);
        outResult.playerB = parsePlayer(result[5]);
        outResult.playerA.position=[parseInt(result[2],10),parseInt(result[3],10),parseInt(result[4],10)];
        outResult.playerB.position=[parseInt(result[6],10),parseInt(result[7],10),parseInt(result[8],10)];
        outResult.weapon = result[9];
        outResult.damage = parseInt(result[10],10);
        outResult.damage_armor = parseInt(result[11],10);
        outResult.health = parseInt(result[12],10);
        outResult.armor = parseInt(result[13],10);
        outResult.hitgroup = result[14];
        return resolve(outResult);
       } else {
       return resolve(undefined);
       };
  });
};

function killedHsFn(line){
  
  return new Promise(function(resolve,reject){
      let result = line.match(/^["](.+)["] \[(.+) (.+) (.+)\] killed ["](.+)["] \[(.+) (.+) (.+)\] with ["](.+)["] (\(headshot\))$/i);

     if (result) {

      let outResult={};
      outResult.event='killed';
      outResult.playerA=parsePlayer(result[1]);
      outResult.playerB=parsePlayer(result[5]);
      outResult.weapon=result[8];
      outResult.playerA.position=[parseInt(result[2],10),parseInt(result[3],10),parseInt(result[4],10)];
      outResult.playerB.position=[parseInt(result[6],10),parseInt(result[7],10),parseInt(result[8],10)];
      outResult.headshot = true;
      outResult.penetrated = false;
      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};
function killedPenHsFn(line){
  
  return new Promise(function(resolve,reject){
      let result = line.match(/^["](.+)["] \[(.+) (.+) (.+)\] killed ["](.+)["] \[(.+) (.+) (.+)\] with ["](.+)["] (\(headshot penetrated\))$/i);

     if (result) {

      let outResult={};
      outResult.event='killed';
      outResult.playerA=parsePlayer(result[1]);
      outResult.playerB=parsePlayer(result[5]);
      outResult.weapon=result[8];
      outResult.playerA.position=[parseInt(result[2],10),parseInt(result[3],10),parseInt(result[4],10)];
      outResult.playerB.position=[parseInt(result[6],10),parseInt(result[7],10),parseInt(result[8],10)];
      outResult.headshot = true;
      outResult.penetrated = true;
      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};
function killedPenFn(line){
  
  return new Promise(function(resolve,reject){
      let result = line.match(/^["](.+)["] \[(.+) (.+) (.+)\] killed ["](.+)["] \[(.+) (.+) (.+)\] with ["](.+)["] (\(penetrated\))$/i);

     if (result) {

      let outResult={};
      outResult.event='killed';
      outResult.playerA=parsePlayer(result[1]);
      outResult.playerB=parsePlayer(result[5]);
      outResult.weapon=result[8];
      outResult.playerA.position=[parseInt(result[2],10),parseInt(result[3],10),parseInt(result[4],10)];
      outResult.playerB.position=[parseInt(result[6],10),parseInt(result[7],10),parseInt(result[8],10)];
      outResult.headshot = false;
      outResult.penetrated = true;
      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};
function killedFn(line){
  
  return new Promise(function(resolve,reject){
    let result = line.match(/^["](.+)["] \[(.+) (.+) (.+)\] killed ["](.+)["] \[(.+) (.+) (.+)\] with ["](.+)["]$/i);
     if (result) {
      let outResult={};
      outResult.event='killed';
      outResult.playerA=parsePlayer(result[1]);
      outResult.playerB=parsePlayer(result[5]);
      outResult.weapon=result[9];
      outResult.playerA.position=[parseInt(result[2],10),parseInt(result[3],10),parseInt(result[4],10)];
      outResult.playerB.position=[parseInt(result[6],10),parseInt(result[7],10),parseInt(result[8],10)];
      outResult.headshot = false;
      outResult.penetrated = false;

      return resolve(outResult);
     } else {
     return  resolve(undefined);
     };
  });
};

function triggerTeamFn(line){
  
  return new Promise(function(resolve,reject){
  let result = line.match(/^Team ["](.+)["] triggered ["](.+)["] [(]CT ["](.+)["][)] [(]T ["](.+)["][)]$/i);
  if (result) {

      let outResult={};
    outResult.event = 'trigger';
    outResult.triggerType = 'team';    
    outResult.trigger = result[2];
    outResult.team = result[1];
    outResult.resultScoreCT = parseInt(result[3],10);
    outResult.resultScoreT = parseInt(result[4],10);

     return  resolve(outResult);
     } else {
     return  resolve(undefined);
     };
  });
};

function triggerWorldOnFn(line){
  
  return new Promise(function(resolve,reject){
      let result = line.match(/^World triggered ["](.+)["] on ["](.+)["]$/i);

  if (result) {
      let outResult={};
    outResult.event='trigger';
    outResult.triggerType ='world';
    outResult.trigger=result[1];
    outResult.triggerOn=result[2];

     return  resolve(outResult);
     } else {
     return  resolve(undefined);
     };
  });
};

function triggerWorldFn(line){
  
  return new Promise(function(resolve,reject){
  let result = line.match(/^World triggered ["](\S+)["]$/i);

  if (result) {
         let outResult={}; 
    outResult.event='trigger';
    outResult.triggerType ='world';
    outResult.trigger = result[1];


      return resolve(outResult);
     } else {
     return  resolve(undefined);
     };
  });
};

function triggerPlayerFn(line){
  
  return new Promise(function(resolve,reject){
  let result = line.match(/^["](.+)["] triggered ["](.+)["]$/i);
  if (result) {
      let outResult={};
    outResult.eventType = 'player';
    outResult.event = 'trigger';
    outResult.player = parsePlayer(result[1]);
    outResult.trigger = result[2];

      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function triggerPlayerValFn(line){
  
  return new Promise(function(resolve,reject){
    
  let result = line.match(/^["](.+)["] triggered ["](.+)["] \(value \"?(.+)\"\)$/i);

  if (result) {

      let outResult={};
    outResult.eventType = 'player';
    outResult.event = 'trigger';
    outResult.player = parsePlayer(result[1]);
    outResult.trigger = result[2];
    outResult.value = result[3];

      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function purchaseFn(line){
  
  return new Promise(function(resolve,reject){
  let result = line.match(/^["](.+)["] purchased ["](.+)["]$/i);

  if (result) {

      let outResult={};
    outResult.event = 'purchase';
    outResult.player = parsePlayer(result[1]);
    outResult.item = result[2];

      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function disconnectedFn(line){
  
  return new Promise(function(resolve,reject){
  let result = line.match(/^["](.+)["] disconnected [(]reason "(.+)"[)]$/i);
  if (result) {
      let outResult={};
  outResult.event='disconnected';
  outResult.player=parsePlayer(result[1]);
  outResult.reason=result[2];

     return  resolve(outResult);
     } else {
     return  resolve(undefined);
     };
  });
};

function enteredFn(line){
  
  return new Promise(function(resolve,reject){
  let result = line.match(/^["](.+)["] entered the game$/i);
  if (result) {
      let outResult={};
  outResult.event='entered';
  outResult.player=parsePlayer(result[1]);

      return resolve(outResult);
     } else {
     return  resolve(undefined);
     };
  });
};

function startLogFn(line){
  
  return new Promise(function(resolve,reject){
   let result = line.match(/^Log file started$/i);

  if (result) {
        let outResult={};  
    outResult.event = 'startLog';

      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function closeLogFn(line){
  
  return new Promise(function(resolve,reject){
  let result = line.match(/^Log file closed$/i);

  if (result) {
      let outResult={};
    outResult.event = 'closeLog';

      return resolve(outResult);
     } else {
     return  resolve(undefined);
     };
  });
};

function connectedFn(line){
  
  return new Promise(function(resolve,reject){
   let result = line.match(/^["](.+)["] connected, address ["]["]$/i);

  if (result) {
      let outResult={};
    outResult.event='connected';
    outResult.player=parsePlayer(result[1]);
    outResult.adress='';

     return  resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function connectedAddressFn(line){
  
  return new Promise(function(resolve,reject){
    let result = line.match(/^["](.+)["] connected, address ["](.+)["]$/i);

  if (result) {
      let outResult={};
    outResult.event='connected';
    outResult.player=parsePlayer(result[1]);
    outResult.adress=result[2];

      return resolve(outResult);
     } else {
     return  resolve(undefined);
     };
  });
};

function rconFn(line){
  
  return new Promise(function(resolve,reject){
  let result = line.match(/^rcon from ["](.+)[:](\d+)["][:] command ["](.+)["]$/i);

  if (result) {
        let outResult={};  
    outResult.event = 'rcon';
    outResult.address = result[1];
    outResult.port = parseInt(result[2],10);
    outResult.command = result[3];

      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function threwFn(line){
  
  return new Promise(function(resolve,reject){
  let result = line.match(/^["](.+)["] threw (.+) \[(.+) (.+) (.+)\]$/i);
  if (result) {

      let outResult={};
    outResult.event = 'threw';
    outResult.player = parsePlayer(result[1]);
    outResult.item = result[2];
    outResult.player.position = [parseInt(result[3],10),parseInt(result[4],10),parseInt(result[5],10)];

      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function pickedUpFn(line){
  
  return new Promise(function(resolve,reject){
  let result = line.match(/^["](.+)["] picked up item ["](.+)["]$/i);
  if (result) {

      let outResult={};
    outResult.event = 'pickedUp';
    outResult.player = parsePlayer(result[1]);
    outResult.item = result[2];

      return resolve(outResult);
     } else {
     return  resolve(undefined);
     };
  });
};

function suicideFn(line){
  
  return new Promise(function(resolve,reject){
    let result = line.match(/^["](.+)["] committed suicide with ["](.+)["]$/i);
  if (result) {

      let outResult={};
    outResult.event = 'suicide';
    outResult.player = parsePlayer(result[1]);
    outResult.item = result[2];

      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function suicidePosition2Fn(line){
  
  return new Promise(function(resolve,reject){
      let result = line.match(/^["](.+)["] \[(.+) (.+) (.+)\] committed suicide with ["](.+)["]$/i);

     if (result) {

      let outResult={};
      outResult.event='suicide';
      outResult.player=parsePlayer(result[1]);
      outResult.player.position=[parseInt(result[2],10),parseInt(result[3],10),parseInt(result[4],10)];
      outResult.item = result[5];

      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function suicidePositionFn(line){
  
  return new Promise(function(resolve,reject){
      let result = line.match(/^["](.+)["] committed suicide with ["](.+)["] [(]attacker_position ["](.+) (.+) (.+)["][)]$/i);
  if (result) {

      let outResult={};
    outResult.event = 'suicide';
    outResult.player = parsePlayer(result[1]);
    outResult.item = result[2];
    outResult.player.position = [parseInt(result[3],10),parseInt(result[4],10),parseInt(result[5],10)];

      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function changedRoleFn(line){
  
  return new Promise(function(resolve,reject){
  let result = line.match(/^["](.+)["] changed role to ["](.+)["]$/i);
  if (result) {

      let outResult={};
    outResult.event = 'changedRole';
    outResult.player = parsePlayer(result[1]);
    outResult.role = result[2];

      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function SIDValidatedFn(line){
  
  return new Promise(function(resolve,reject){
    let result = line.match(/^["](.+)["] STEAM USERID validated$/i);
  if (result) {

      let outResult={};
    outResult.event = 'SIDValidated';
    outResult.player = parsePlayer(result[1]);

      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function joinedTeamFn(line){
  
  return new Promise(function(resolve,reject){
  let result = line.match(/^["](.+)["] joined team ["](.+)["]$/i);
  if (result) {

      let outResult={};
    outResult.event = 'joinedTeam';
    outResult.player = parsePlayer(result[1]);
    outResult.team = result[2];

      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function switchTeamFn(line){
  
  return new Promise(function(resolve,reject){
  let result = line.match(/^["](.+)["] switched from team [<](.+)[>] to [<](.+)[>]$/i);
  if (result) {

      let outResult={};
    outResult.event = 'switchTeam';
    outResult.player = parsePlayer(result[1]);
    outResult.teamA = result[2];
    outResult.teamB = result[3];


      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function loadingMapFn(line){
  
  return new Promise(function(resolve,reject){
      let result = line.match(/^Loading map ["](.+)["]$/i);
    if (result) {
  

      let outResult={};
    outResult.event = 'loadingMap';
    outResult.map = result[1];
    // console.log(outResult);
      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function startedMapFn(line){
  
  return new Promise(function(resolve,reject){
  let result = line.match(/^Started map ["](.+)["] [(]CRC ["](.+)["][)]$/i);
  
  if (result) {
      let outResult={};
    outResult.event = 'startedMap';
    outResult.map = result[1];
    outResult.crc = result[2];
    
    // console.log(outResult);

      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function serverCvarFn(line){
  
  return new Promise(function(resolve,reject){
    let result = line.match(/^server_cvar: ["](.+)["] ["](.+)["]$/i);
  if (result) {

      let outResult={};
    outResult.event = 'serverCvar';
    outResult.parameter = result[1];
    outResult.value = result[2];

      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function scoreCurrentFn(line){
  
  return new Promise(function(resolve,reject){
  let result = line.match(/^Team ["](.+)["] current score ["](\d+)["] with ["](\d+)["] players$/i);
  if (result) {

      let outResult={};
    outResult.event = 'score';
    outResult.team = result[1];
    outResult.scoreType = 'current'
    outResult.score = parseInt(result[2],10);
    outResult.players = parseInt(result[3],10);

      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function scoreScoredFn(line){
  
  return new Promise(function(resolve,reject){
  let result = line.match(/^Team ["](.+)["] scored ["](\d+)["] with ["](\d+)["] players$/i);
  if (result) {

      let outResult={};
    outResult.event = 'scored';
    outResult.team = result[1];
    outResult.scoreType = 'scored'
    outResult.score = parseInt(result[2],10);
    outResult.players = parseInt(result[3],10);

     return  resolve(outResult);
     } else {
     return  resolve(undefined);
     };
  });
};

function scoreFinalFn(line){
  
  return new Promise(function(resolve,reject){
    let result = line.match(/^Team ["](.+)["] final score ["](\d+)["] with ["](\d+)["] players$/i);
  if (result) {

      let outResult={};
    outResult.event = 'scored';
    outResult.team = result[1];
    outResult.scoreType = 'final'
    outResult.score = parseInt(result[2],10);
    outResult.players = parseInt(result[3],10);

     return  resolve(outResult);
     } else {
    return   resolve(undefined);
     };
  });
};

function spawnedFn(line){
  
  return new Promise(function(resolve,reject){
  let result = line.match(/^["](.+)["] spawned as ["](.+)["]$/i);
  if (result) {

      let outResult={};
    outResult.event = 'spawned';
    outResult.player = parsePlayer(result[1]);
    outResult.spawned = result[2];

     return  resolve(outResult);
     } else {
     return  resolve(undefined);
     };
  });
};

function sayFn(line){
  
  return new Promise(function(resolve,reject){
  let result = line.match(/^["](.+)["] say ["](.+)["]$/i);
  if (result) {

      let outResult={};
    outResult.event = 'say';
    outResult.player = parsePlayer(result[1]);
    outResult.text = result[2];

     return  resolve(outResult);
     } else {
     return  resolve(undefined);
     };
  });
};

function sayTeamFn(line){
  
  return new Promise(function(resolve,reject){
 
  let result = line.match(/^["](.+)["] say_team ["](.+)["]$/i);
  if (result) {

      let outResult={};
    outResult.event = 'sayTeam';
    outResult.player = parsePlayer(result[1]);
    outResult.text = result[2];

     return  resolve(outResult);
     } else {
     return  resolve(undefined);
     };
  });
};

function positionReportFn(line){
  
  return new Promise(function(resolve,reject){
    let result = line.match(/^["](.+)["] position_report (.+)$/i);
  if (result) {

      let outResult={};
    outResult.event = 'positionReport';
    outResult.player = parsePlayer(result[1]);
    outResult.text = result[2];

     return  resolve(outResult);
     } else {
     return  resolve(undefined);
     };
  });
};

function assistFn(line){
  
  return new Promise(function(resolve,reject){
    let result = line.match(/^["](.+)["] assisted killing ["](.+)["]$/i);
  if (result) {

      let outResult={};
    outResult.event = 'assist';
    outResult.playerA = parsePlayer(result[1]);
    outResult.playerB = parsePlayer(result[2]);

      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function voteSuccessKickFn(line){
  
  return new Promise(function(resolve,reject){
    let result = line.match(/^Vote succeeded ["]Kick (.+)["]$/i);
  if (result) {

    let outResult={};
    outResult.event = 'voting';
    outResult.voteType = 'kick';    
    outResult.voteResult = 'success';
    outResult.player = result[1];


      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};
function voteFalseKickFn(line){
  
  return new Promise(function(resolve,reject){
    let result = line.match(/^Vote failed ["]Kick (.+)["]$/i);
  if (result) {

    let outResult={};
    outResult.event = 'voting';
    outResult.voteType = 'kick';        
    outResult.voteResult = 'failed ';
    outResult.player = result[1];


      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function voteSuccessNxLvlFn(line){
  
  return new Promise(function(resolve,reject){
    let result = line.match(/^Vote succeeded ["]NextLevel (.+)["]$/i);
  if (result) {

    let outResult={};
    outResult.event = 'voting';
    outResult.voteType = 'nextlevel';    
    outResult.voteResult = 'success';
    outResult.map = result[1];


      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function voteFalseNxLvlFn(line){
  
  return new Promise(function(resolve,reject){
    let result = line.match(/^Vote failed ["]NextLevel (.+)["]$/i);
  if (result) {

    let outResult={};
    outResult.event = 'voting';
    outResult.voteType = 'nextlevel';        
    outResult.voteResult = 'failed ';
    outResult.map = result[1];


      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function voteSuccessChLvlFn(line){
  
  return new Promise(function(resolve,reject){
    let result = line.match(/^Vote succeeded ["]ChangeLevel (.+)["]$/i);
  if (result) {

    let outResult={};
    outResult.event = 'voting';
    outResult.voteType = 'changelevel';    
    outResult.voteResult = 'success';
    outResult.map = result[1];


      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};
function voteFalseChLvlFn(line){
  
  return new Promise(function(resolve,reject){
    let result = line.match(/^Vote failed ["]ChangeLevel (.+)["]$/i);
  if (result) {

    let outResult={};
    outResult.event = 'voting';
    outResult.voteType = 'changelevel';        
    outResult.voteResult = 'failed ';
    outResult.map = result[1];


      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};
function voteSuccessScrambleFn(line){
  
  return new Promise(function(resolve,reject){
    let result = line.match(/^Vote succeeded ["]ScrambleTeams ["]$/i);
  if (result) {

    let outResult={};
    outResult.event = 'voting';
    outResult.voteType = 'scrambleteams';    
    outResult.voteResult = 'success';

      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};
function voteFalseScrambleFn(line){
  
  return new Promise(function(resolve,reject){
    let result = line.match(/^Vote failed ["]ScrambleTeams ["]$/i);
  if (result) {

    let outResult={};
    outResult.event = 'voting';
    outResult.voteType = 'scrambleteams';        
    outResult.voteResult = 'failed ';

      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function voteSuccessSwapFn(line){
  
  return new Promise(function(resolve,reject){
    let result = line.match(/^Vote succeeded ["]SwapTeams ["]$/i);
  if (result) {

    let outResult={};
    outResult.event = 'voting';
    outResult.voteType = 'swapteams';    
    outResult.voteResult = 'success';

      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};
function voteFalseSwapFn(line){
  
  return new Promise(function(resolve,reject){
    let result = line.match(/^Vote failed ["]SwapTeams ["]$/i);
  if (result) {

    let outResult={};
    outResult.event = 'voting';
    outResult.voteType = 'swapteams';        
    outResult.voteResult = 'failed ';

      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};

function molotovFn(line){
  
  return new Promise(function(resolve,reject){
    let result = line.match(/^Molotov projectile spawned at (.+) (.+) (.+), velocity (.+) (.+) (.+)$/i);
  if (result) {

      let outResult={};
    outResult.event = 'molotovProjectile';
    outResult.position = [parseInt(result[1],10),parseInt(result[2],10),parseInt(result[3],10)];
    outResult.positionVelocity = [parseInt(result[4],10),parseInt(result[5],10),parseInt(result[6],10)];

      return resolve(outResult);
     } else {
      return resolve(undefined);
     };
  });
};


/////////////////////////// hlx2
function triggerKillLocationFn(line){
  
  return new Promise(function(resolve,reject){
  let result = line.match(/^World triggered ["](.+)["] \(attacker_position ["](.+) (.+) (.+)["]\) \(victim_position ["](.+) (.+) (.+)["]\)$/i);

  if (result) {
    let outResult={}; 
    
    outResult.event='trigger';
    outResult.triggerType ='world';
    outResult.trigger = result[1];
    outResult.playerA.position = [parseInt(result[2],10),parseInt(result[3],10),parseInt(result[4],10)];
    outResult.playerB.position = [parseInt(result[5],10),parseInt(result[6],10),parseInt(result[7],10)];

      return resolve(outResult);
     } else {
     return  resolve(undefined);
     };
  });
};

function triggerWeaponstatsFn(line){
  
  return new Promise(function(resolve,reject){
    let result = line.match(/^["](.+)["] triggered ["](.+)["] \(weapon ["](.+)["]\) \(shots ["](.+)["]\) \(hits ["](.+)["]\) \(kills ["](.+)["]\) \(headshots ["](.+)["]\) \(tks ["](.+)["]\) \(damage ["](.+)["]\) \(deaths ["](.+)["]\)$/i);

    if (result) {
      let outResult={}; 
      outResult.event='trigger';
      outResult.triggerType ='Weaponstats';
      outResult.playerA = parsePlayer(result[1]);
      outResult.weapon =  result[3];
      outResult.shots =  result[4];
      outResult.hits =  result[5];
      outResult.kills =  result[6];
      outResult.headshots =  result[7];
      outResult.tks =  result[8];
      outResult.damage =  result[9];    
      outResult.deaths =  result[10];    

      return resolve(outResult);
    } else {
      return  resolve(undefined);
    };

  });
};

function triggerWeaponstats2Fn(line) {


  return new Promise((resolve,reject)=>{
    let result = line.match(/^["](.+)["] triggered ["](.+)["] \(weapon ["](.+)["]\) \(head ["](.+)["]\) \(chest ["](.+)["]\) \(stomach ["](.+)["]\) \(leftarm ["](.+)["]\) \(rightarm ["](.+)["]\) \(leftleg ["](.+)["]\) \(rightleg ["](.+)["]\)$/i);

    if (result) {
      let outResult={}; 
      outResult.event='trigger';
      outResult.triggerType ='Weaponstats2';
      outResult.playerA = parsePlayer(result[1]);
      outResult.weapon =  result[3];
      outResult.head =  result[4];
      outResult.chest =  result[5];
      outResult.stomach =  result[6];
      outResult.leftarm =  result[7];
      outResult.rightarm =  result[8];
      outResult.leftleg =  result[9];    
      outResult.rigthleg =  result[10];    

      return resolve(outResult);
    } else {
      return  resolve(undefined);
    };
  });
}

///////////////////////////

/// main parser
module.exports = (line) => {
  


 if (Buffer.isBuffer(line) === true) {
    let line = line.toString('utf8');
  } 

let idx = line.indexOf('L ');

if (idx == -1) return undefined;

let line = line.substring(30).replace(/(\r\n|\n|\r|\u0000)/gm,"").trim();

//////////////////////////////////////////////////////

return Promise.all([
  attakedFn(line),
  killedHsFn(line),
  killedFn(line),
  killedPenHsFn(line),
  killedPenFn(line),
  molotovFn(line),
  teamNameFn(line),
  voteSuccessKickFn(line),
  voteFalseKickFn(line),
  voteSuccessNxLvlFn(line),
  voteFalseNxLvlFn(line),
  voteSuccessChLvlFn(line),
  voteFalseChLvlFn(line),
  voteSuccessScrambleFn(line),
  voteFalseScrambleFn(line),
  voteSuccessSwapFn(line),
  voteFalseSwapFn(line),
  triggerTeamFn(line),
  triggerWorldOnFn(line),
  triggerWorldFn(line),
  triggerPlayerFn(line),
  triggerPlayerValFn(line),
  purchaseFn(line),
  disconnectedFn(line),
  enteredFn(line),
  startLogFn(line),
  closeLogFn(line),
  connectedFn(line),
  connectedAddressFn(line),
  rconFn(line),
  threwFn(line),
  pickedUpFn(line),
  suicideFn(line),
  suicidePosition2Fn(line),
  suicidePositionFn(line),
  changedRoleFn(line),
  SIDValidatedFn(line),
  joinedTeamFn(line),
  switchTeamFn(line),
  loadingMapFn(line),
  startedMapFn(line),
  serverCvarFn(line),
  scoreCurrentFn(line),
  scoreScoredFn(line),
  scoreFinalFn(line),
  spawnedFn(line),
  sayFn(line),
  sayTeamFn(line),
  positionReportFn(line),
  assistFn(line),
  triggerKillLocationFn(line),
  triggerWeaponstatsFn(line),
  triggerWeaponstats2Fn(line)
  ])
  .then((res) => {
/// filtering undefined strings
    let msg = res.filter( (res) => res );
////////////////////////////

    if (msg.length > 1 ){

//////////// loging strings that wrong parsed
      fs.writeFile('./logs/wrong_parsing.txt',JSON.stringify(msg)+'\n',{flag:'a'});
      fs.writeFile('./logs/wrong_parsing.txt',JSON.stringify(line)+'\n',{flag:'a'});      
////////////////////////////

    } else if (msg.length === 1 ) {

      let outResult={};

      let d = new Date();
      let year = d.getFullYear();
      let month = d.getMonth();
      let day = d.getDay();
      let hours = d.getHours();
      let minutes = d.getMinutes();
      let seconds = d.getSeconds();

      outResult = msg[0];
      outResult.date = {
        year : year,
        month : month,
        day : day,
        hour : hours,
        minutes : minutes,
        seconds : seconds,
        ts : +d,
        fulldate : d
    };
/// retrun parsed JSON promise
      return outResult;

    } else {

//////////// loging strings that not parsed
      fs.writeFile('./logs/not_parsed.txt',JSON.stringify(line)+'\n',{flag:'a'});
///////////////////////

    }
  })
  .catch( (err) => console.log(err) );
}


function parsePlayer(line) {

    let result = line.match(/^(.+)[<](\d+)[>][<]([STEAM1234567890_:]+)[>][<](.+)[>]$/);

    if (result) { 
      return {
        name : result[1],
        id : parseInt(result[2],10),
        steamid : result[3],
        team : result[4]
        }
	 };

    result = line.match(/^(.+)[<](\d+)[>][<]([STEAM1234567890_:]+)[>][<][>]$/);
    
    if (result) { 
      return {
        name : result[1],
        id : parseInt(result[2],10),
        steamid : result[3],
        team : 'Unassigned'
        }
	};

    result = line.match(/^(.+)[<](\d+)[>][<]([STEAM1234567890_:]+)[>]$/);

    if (result) { 
      return {
        name : result[1],
        id : parseInt(result[2],10),
        steamid : result[3],
        team : 'Unassigned'
        }
	};
    result = line.match(/^(.+)[<](\d+)[>][<]([BOT]+)[>][<](.+)[>]$/);

    if (result) { 
      return {
        name : result[1],
        id : parseInt(result[2],10),
        steamid : result[3],
        team : result[4]
        }
   };

    result = line.match(/^(.+)[<](\d+)[>][<]([BOT]+)[>][<][>]$/);
    
    if (result) { 
      return {
        name : result[1],
        id : parseInt(result[2],10),
        steamid : result[3],
        team : 'Unassigned'
        }
  };

    result = line.match(/^(.+)[<](\d+)[>][<]([BOT]+)[>]$/);

    if (result) { 
      return {
        name : result[1],
        id : parseInt(result[2],10),
        steamid : result[3],
        team : 'Unassigned'
        }
  };
  
    result = line.match(/^(.+)[<](\d+)[>][<]([consoleCONSOLE]+)[>]$/);

    if (result) { 
      return {
        name : result[1],
        id : parseInt(result[2],10),
        steamid : result[3],
        team : 'Unassigned'
        }
    };
  
    result = line.match(/^(.+)[<](\d+)[>][<]([consoleCONSOLE]+)[>][<]([consoleCONSOLE]+)[>]$/);

    if (result) { 
      return {
        name : result[1],
        id : parseInt(result[2],10),
        steamid : result[3],
        team : result[4]
        }
    };  
  return undefined;

};
/*
function parseIP(line){

	let result = line.match(/(.+):(.+)/);

	if (result) {
		return {'ip':result[1],'port':result[2]}; 
	} else {
        return undefined;
    }
};

function parseArgs(args) {
  
  let result = args.match(/[(]([^"]|\\")+ ["]([^"]|\\")+["][)]/g);
  
  if (result) {
        return result;
    } else {
        return undefined;
    };

};

*/