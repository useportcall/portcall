#!/bin/sh
set -e

npx prisma init

# Pull the latest schema from the database
npx prisma db pull

# Start Prisma Studio
exec npx prisma studio --port 5555
