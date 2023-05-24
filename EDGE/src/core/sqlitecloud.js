const logThis = (id = "SQLiteCloud", msg) => {
  const dateObject = new Date();
  // adjust 0 before single digit date
  const date = (`0 ${dateObject.getDate()}`).slice(-2);
  // current month
  const month = (`0${dateObject.getMonth() + 1}`).slice(-2);
  // current year
  const year = dateObject.getFullYear();
  // current hours
  const hours = dateObject.getHours();
  // current minutes
  const minutes = dateObject.getMinutes();
  // current seconds
  const seconds = dateObject.getSeconds();
  // current milliseconds
  const milliseconds = dateObject.getMilliseconds()
  // prints date & time in YYYY/MM/DD HH:MM:SS format
  const prefix = `${year}/${month}/${date} ${hours}:${minutes}:${seconds}:${milliseconds}`;
  console.log(`!!!!!!!!! ${id}: ${prefix} - ${msg}`);
}
/*
SQLiteCloud class
*/
export default class SQLiteCloud {
  /* PRIVATE PROPERTIES */
  /* 
  #debug_sdk 
  */
  #debug_sdk = false;
  /* CONSTRUCTOR */
  /*
  SQLiteCloud class constructor receives:
   - project ID (required)
   - api key (required)
   - debug
  */
  constructor(projectID, apikey, debug_sdk = false) {
    this.#debug_sdk = debug_sdk;
    this.url = `wss://web1.sqlitecloud.io:8443/api/v1/${projectID}/ws?apikey=${apikey}`;
  }
  /*
  method used to send SQL commands using fetch API
  */
  async sendCommands(commands) {
    if (this.#debug_sdk) logThis("", "sending commands: " + commands);
    const data = {
      "statement": commands
    }
    const response = await fetch(this.url, {
      method: "POST",
      mode: "cors",
      cache: "no-cache",
      credentials: "same-origin",
      headers: {
        "Content-Type": "application/json",
      },
      redirect: "follow",
      referrerPolicy: "no-referrer",
      body: JSON.stringify(data),
    });
    const jsonData = await response.json();
    return jsonData
  }
}