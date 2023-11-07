FROM mongo:5.0.22
# Will be set through Environment Files
ARG HOST 
ARG PORT
ARG DATABASE
ARG DATABASE_URI

ENV HOST ${HOST}
ENV PORT ${PORT}
ENV DATABASE ${DATABASE}
ENV DATABASE_URI ${DATABASE_URI}

COPY user.json /user.json
CMD mongoimport --uri ${DATABASE_URI} --db ${DATABASE} --collection user --type json --file /user.json --jsonArray