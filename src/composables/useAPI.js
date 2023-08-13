export function useApiUrl(path,ws) {
  const currentHost = window.location.hostname;
  const subdomain = ()=>{
    if(window.location.subdomain){
      return window.location.subdomain
    }
    return ''
  };
  const domain = currentHost.replace(/^www\./, ''); // Remove 'www.' if present
  let protocol = window.location.protocol;
  if (ws){
    if(domain === 'localhost'){
      protocol = 'ws:'
    }
    else{
      protocol = 'wss:'
    }
  }

  return `${protocol}//${subdomain()}${domain}/api/v1${path}`;
};
