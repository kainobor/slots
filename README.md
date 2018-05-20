# SLOTS
Simple REST-API for slot game

# Task description

It should be REST api with 1 implemented action to create SPIN

POST /api/machines/atkins-diet/spins

Request is JWT token with this payload:
```json
{
	uid: “asfasf”,  // user id
	chips: 10000, // chips balance
	bet: 1000 // bet size
}
```

Result:
```json
{
    "total": 0, - total win
    "spins": [
        {
            "type": "main", // can be main or free spins
            "total": 0, // winning in spin
            "stops": [ // real stops
                13,
                18,
                9
            ]
        }
    ],
    "jwt": {
    	uid: “asfasf”,  // user id
	    chips: 10000,
	    bet: 1000
    }
}
```

On each request service should take user balance in chips, decrease it by bet size and calculate spin or spins (if user got free spins) result. And generate new jwt token with new balance for response.

Service should be designed extendable to support different machines (currently it’s only one)
Service should be implemented with high load in mind.

To get paylines please mouse over the numbers on the left or right here https://wizardofodds.com/play/slots/atkins-diet/  You will see how paylines are defined.
You can actually play this machine and test payout with your implementation (just press play).