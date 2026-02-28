# Use Node image for development
FROM node:20-alpine

WORKDIR /app

# Copy package files and install dependencies
COPY package*.json ./
RUN npm ci

# Copy the rest of the source code (will be overridden by volume mount)
COPY . .

# Expose the Vite dev server port (default 5173)
EXPOSE 5173

# Start the development server with hot reload
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]