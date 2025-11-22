#!/bin/sh
set -e

# if prisma folder does not exist, initialize it
[ ! -d "./prisma" ] && npx prisma init

# Remove 'url' from schema.prisma if present (Prisma 7+)
if [ -f "./prisma/schema.prisma" ]; then
  sed -i '/url *=/d' ./prisma/schema.prisma
fi

# Create prisma.config.ts if it doesn't exist
if [ ! -f "./prisma/prisma.config.ts" ]; then
  cat <<EOF > ./prisma/prisma.config.ts
import { defineConfig } from '@prisma/internals';
export default defineConfig({
  datasource: {
    provider: 'postgresql',
    url: process.env.DATABASE_URL,
  },
});
EOF
fi

# log something
echo "Starting Prisma Studio..."

# Pull the latest schema from the database
prisma db pull

# Start Prisma Studio
exec prisma studio --port 5555 --browser none