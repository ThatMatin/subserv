Subscription service
The task is to build a service including the features of the following stories.
You can choose the programing language you want, for us it is important to see how you work and what approaches you have. The data can be
mocked, you can use any kind of database you want to use.
Please use a git repository to keep track on your changes and share the link with us when you have everything finished.
We are expecting an API that matches the following user behaviours. The user API is called by our client e.g. OTT Apps or Web platform (that you
don’t have to build).
Please create a documentation that we can see all created API endpoints and how we are able to use them. Also share a README.md with us,
how we can test your service..

User Story 1
As a User, I want to be able to select from a list a product, and based on this product to receive a subscription plan.
AC:
1. I can fetch a list of products
2. I can fetch a single product
3. I can buy a single product
4. I want to fetch the following informations related to my subscription (e.g. start date, end date, duration of the subscription, prices, tax)
5. I can pause and unpause my subscription
6. I can cancel my active subscription
Technical details:
- Each product from the list should have a different subscription duration and a price to reflect the length of the subscription duration

Optional Story 1
As a User, I want to be able to use a voucher code to buy products with a discounted price (fixed amount and percentage discount)
AC:
1. I can list products with individual voucher
2. I can buy a product with a voucher
3. A validation for the voucher should be in place
4. Subscription plan is created with the discounted price
5. I want to fetch the following informations related to my subscription (e.g. start date, end date, duration of the subscription, prices, tax)

Optional Story 2
As a User, I want to receive a trial period of 1 month before the start of my subscription
AC:
1. I want to fetch the following informations related to my subscription (e.g. start date, end date, duration of the subscription, prices, tax,
trial)
2. During the trial I can cancel my subscription
3. During the trial I am not able to pause my subscription
