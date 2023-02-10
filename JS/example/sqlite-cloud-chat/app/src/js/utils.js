import { format } from 'date-fns';
/*
debug utility 
*/
const logThis = (msg) => {
  let prefix = format(new Date(), "yyyy-MM-dd'T'HH:mm:ss.SSS");
  console.log(prefix + " - " + msg);
}
export { logThis };
/*
build css classes removing null and undefined values 
*/
const classNames = (...classes) => {
  return classes.filter(Boolean).join(' ')
}
export { classNames };
/*
this method is used to check if a channelName (e.g. a channel indicated in a queryString) exists in the list of avaible channels
if channel exists, return the index of channelName in the array channelsList
if channel not exists, return -1
*/
const checkChannelExist = (channelsList, channelName) => {
  let chIndex = -1;
  if (channelsList) {
    channelsList.forEach((ch, i) => {
      if (ch == channelName) chIndex = i;
    })
    return chIndex;
  } else {
    return chIndex;
  }
}
export { checkChannelExist };

