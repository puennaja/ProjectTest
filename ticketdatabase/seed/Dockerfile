FROM mongo:5.0.22
# Will be set through Environment Files
ARG DATABASE_URI

ENV DATABASE_URI ${DATABASE_URI}

COPY user.json /user.json
CMD mongoimport --uri ${DATABASE_URI} --collection user --drop --file /user.json