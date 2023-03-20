## Background
You have been tasked with building a reverse proxy that will help with data governance for our
startup. Your task is to create a simple HTTP server that will act as a reverse proxy.
The proxy should be able to handle incoming requests, inspect the request headers, and then
forward the request to the intended destination. In addition, the proxy should be able to log all
incoming requests, including the headers and body, and the response headers and body. The
proxy should also be able to block certain requests based on a set of predefined rules.

## Requirements:

1. The reverse proxy should be able to forward any kind of request but inspect only GET
requests.
2. The proxy should inspect the request headers for any sensitive information and mask it
before forwarding the request.
3. The proxy should log all incoming requests, including headers and body, and response
headers and body.
4. The proxy should be able to block requests based on a set of predefined rules. For
example, you can block requests that contain specific headers or parameters.
5. You should provide a configuration file that allows the startup to define the rules for
blocking requests.
6. You should implement the reverse proxy in Golgang or Rust. You are free to use any
libraries or frameworks you see fit.
7. You should include unit tests to ensure that the proxy is functioning correctly.
8. You should also provide documentation on how to install and run the reverse proxy.

## Evaluation Criteria:

1. The code should be clean, readable, and maintainable.
2. The proxy should handle requests quickly and efficiently.
3. The proxy should log all incoming requests, including headers and body, and response
headers and body.
4. The proxy should be able to block requests based on predefined rules.
5. The configuration file should be easy to understand and modify.
6. The unit tests should cover all aspects of the proxy's functionality.
7. The documentation should be clear and concise.
