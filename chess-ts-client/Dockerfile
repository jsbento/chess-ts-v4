FROM node:22.20.0

WORKDIR /app

# Install curl for healthchecks
RUN apt-get update && apt-get install -y curl && rm -rf /var/lib/apt/lists/*

# Copy package files
COPY package.json package-lock.json ./

# Install dependencies
RUN npm install

COPY . .

EXPOSE 5173

# For development, use dev command with hot reloading
CMD ["npm", "run", "dev"]