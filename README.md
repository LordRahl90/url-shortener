# Coding Challenge: URL Shortener in Go

### What is the mission?
Our marketing partners want to show short, easy to type URLs on the info screen in our vehicles.

### What are you supposed to do?
Design and implement an URL shortener HTTP service that fulfills the following criteria:

* Provides an HTTP API to:
	* Shorten a URL
	* Redirect to the long URL from the shortened URL
* Shortened URL requirements:
	* The ID of the shortened URL needs to be unique (across past and concurrent requests)
	* The ID of the shortened URL should be as short as possible (max. 10 characters long)

For example when your service would run under the domain https://short.io this URL https://www.moia.io/my/ride/123456 should look like the this https://short.io/a2bj3 after being shortened. When calling https://short.io/a2bj3 a user should be redirected to https://www.moia.io/my/ride/123456.

### Acceptance Criteria
* Any URL can be sent to the REST service which will return a shortened URL.
* The shortened URL has a length of at most 10 characters.
* A shortened URL cannot be duplicated; e.g. two different links won’t provide the same short name.
* Your solution is implemented in Go.
* Persisting the shortened Urls is optional. It is also enough to just store it in-memory.
* We have to be able to run your solution during assessment.
* Authentication or authorization can be ignored here.

### Supporting information
* What we look at during the assessment:
	* How easy it is to understand and maintain your code.
	* How you verify your software, whether by automated tests or otherwise.
	* How clean your design and implementation is.
* Document where you think it makes sense and helps us to review your code.
