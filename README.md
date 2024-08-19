# Auction

Develop an algorithm for a computerized auction site where sellers offer items for sale, and buyers bid against each other. The goal is to automatically determine the winning bid after all bidders have entered their information. The API will be integrated into the auction site by other developers.

## Requirements

### Starting Bid

- The initial and lowest bid a buyer is willing to offer for the item.

### Max Bid

- The maximum amount the bidder is willing to pay for the item.

### Auto-increment Amount

- A specified dollar amount that the algorithm will add to the bidder's current bid each time they are outbid.
The algorithm should not exceed the bidder's Max Bid.
Increments must be exactly equal to the auto-increment amount.

## Algorithm Goals

### Winning Bid Calculation

- Determine the winning bid as the lowest possible amount that still wins the auction, following all specified rules.
Tiebreaker

- Resolve ties by giving priority to the earliest bid received.

By adhering to these requirements, the algorithm will ensure a fair and efficient bidding process, optimizing for the lowest winning bid while respecting each bidder's constraints.

## How to run
To run this project using Docker, follow these steps:

1. Make sure you have Docker installed on your machine.
2. Clone this repository and navigate to root of this reporsitory using your terminal.
3. Build the Docker image by running the following command:
    ```
    docker build -t auction .
    ```
4. Once the image is built, you can run the project using the following command:
    ```
    docker run -p 8080:8080 auction
    ```
    This will start the project and map port 8080 of the container to port 8080 of your local machine.
5. Open your web browser and visit `https://editor.swagger.io/`, copy all contents from `openapi.yml` to the editor to access the project.

> Note: You might need to disable CORS on your browser for this to work. And, if you're using chrome you can try this extension - https://chromewebstore.google.com/detail/allow-cors-access-control/lhobafahddgcelffkeicbaginigeejlf
