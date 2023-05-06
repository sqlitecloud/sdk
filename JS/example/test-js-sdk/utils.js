// import { format } from 'date-fns';
/*
debug utility 
*/
const logThis = (msg) => {
  let prefix = window.moment().format('YYYY/MM/DD HH:mm:ss.SSS');
  console.log(prefix + " - " + msg);
}
export { logThis };