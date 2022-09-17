addEventListener('fetch', event => {
  event.respondWith(handleRequest(event.request))
})

async function handleRequest(request) {
  // new URL object to play with,
  // based on the one being requested.
  // e.g. https://domain.com/blog/page
  var url = new URL(request.url)
  // set hostname to the place we're proxying requests from
  url.hostname = "${hostname}"

  // remove the first occurence of /blog
  // so it requests / of the proxy domain
  //url.pathname = url.pathname.replace('/blog', '')
  // pass the modified url back to the request,
  let response = await fetch(url, request)
  return response;
}
