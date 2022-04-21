# Poke API service description



## General

Develop a REST service that allows a trainer to insure only their fire, water, and grass-type Pokémon.

To use the service, the trainer must enter their name, surname, and trainer code, a unique and secret code issued by the Pokémon League. This API will return a random and unique MT code to be used to authorize all subsequent API calls.

To get a quote, the trainer will have to enter the name of the Pokémon to be insured, one at a time. The trainer will have to pay for each Pokémon insured an amount equal to € 9.50, plus € 0.50 if the Pokémon should also be of the flying type. A quote will therefore be identified by the price for each insured Pokémon. If the Pokémon isn't insurable due to its type, return a specific error.

Finally, the web service must allow the trainer to purchase the quote and record the effective transaction date. Integration with any payment system is not required for this test.

## Project info

The project has been developed in Go.
The Database is locally, created runtime, based on sql-lite.
The Pokemon API used to retrieve the pokemon information is reached using a 3PP go client (https://github.com/mtslzr/pokeapi-go).

### Project Structure

**cmd:** contains the file, main.go, that handle the entry point of the service. It's responsible to initialize the structure used by the service, and start the server.

 **config:** contains the file, config.go, used for the behaviour based on ENV variable. From the input specific, it's used to mock the payment system based on a dedicated ENV variable: PAYMENT_BEHAVIOUR. The default value of the variable is "accept", which implies the acceptance of the payment. If the variable is set to "rejected", the payment system will return a failure, so it will not be possible to pay an existing quote.

 **db:** contains the files used to handle all the operations related to the DB. It's responsible to initialize the structure related to the connection and expose the interface to interact with the DB, in order to performs all SQL operations.

 **handlerMessage:** contains the file used to handle all the errors related to the service logic. Each errors is based on a specific case, and it's coupled with a dedicated error message and http status code.

 **payment:** contains the file used to handle the payment system. It's exposed an interface to request a payment.

**poke:** contains the file used to interact with Poke API and for the validation of pokemon category enabled to be insured.

**server:** contains the file to create the server and handle the client requests.


### Implementation choices

The project is completely based on the logic of the interfaces that each file of a folder exposes to files of other folders. Each folder, as indicated in the structuring of the project, is responsible for an entity of the service and all its logic, requiring externally managed functionalities through interfaces.
This design choice, together with the use of data structures that in turn contained interfaces as properties, makes it easier to reuse the code and add future changes.

A strong point of dynamism is also represented by the choice not to fix the pokemon quotes in the code, but to keep them in an appropriate table (see paragraph XX), so that future price variations do not change the logic of the methods for the calculation of quotes. As far as the categories are concerned, it was decided to leave them in a list dedicated and initialized to the start-up, but the methods that calculate the admissible categories refer only to a map initialized by an external list, therefore without having linked a specific list of categories to the application logic of the method. The use of the map thus allows a search in O (1) time.

For simplification purposes, it was decided to generate the unique authentication mt as the result of (idTrainer + 1), with idTrainer representing the identifier of the trainer provided during registration. The mt is expected to be put in header client request.

Downstream of the requirements, it seemed reasonable to assume that it was superfluous to allow you to quote a pokemon that has already received a quote. Starting from the requirements, it was not considered appropriate to implement a logic in which a new quotation was accepted only if the price had changed. Different speech for the insurance, so starting from the assumption that in the real context it is possible to insure an entity several times, the possibility of reassuring the pokemon is not forbidden, especially considering that each transition is linked to a timestamp, which leaves open the possibility in the future to think about a deadline.

### Database Model

**Trainer**

Table to store all the trainer.


* id_trainer **unique id provided by the trainer(Primary Key)**
* name: **name of the trainer.**
* surname: **surname of the trainer.**
* mt: **unique authorization code used to interact with the API after the registration step of the trainer.**

**Quote**

Table to store all the requested quotes.

* id_trainer **unique id provided by the trainer**
* pokemon: **name of the pokemon quoted**
* price: **total price of the quote**

(id_trainer,pokemon): PRIMARY KEY
id_trainer: FOREIGN KEY, referred to table TRAINER(id_trainer)

**PayedQuote**

Table to store all the payed quotes.

* id_transaction **unique id of the transaction**
* id_trainer **unique id provided by the trainer**
* pokemon: **name of the pokemon quoted**
* price: **total price of the quote**
* timestamp: **event time to record the transaction**

id_transaction: PRIMARY KEY
id_trainer: FOREIGN KEY, referred to table TRAINER(id_trainer)

**BaseQuote**

Table to store all the payed quotes for mandatory categories for insurance purpose.

* category: **category pokemon**
* price: **price related to category**

(category, price): PRIMARY KEY

**ExtraQuote**

Table to store all the payed quotes for optional categories which lead to an extra price.

* category: **category pokemon**
* price: **price related to category**

(category, price): PRIMARY KEY
 

### Automated test

Unit tests were carried out on all functions, and for each one the possible inputs were considered that would allow to cover all the computation branches.
The use of mocks was made in tests where external functions were interrogated by the responsibility of the method under analysis.
All unit tests are present in their respective folders, while a complete function test is present in the function_test folder.

### Local start

go run ./cmd

### Docker start

docker build -t <image_name> .

docker run -p 8080:8080 -t <image_name> .