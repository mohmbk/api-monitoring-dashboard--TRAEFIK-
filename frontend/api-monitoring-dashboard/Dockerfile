
FROM node:20-alpine AS builder
WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

# Copie le reste du code et build
COPY . .
RUN npm run build 
# la on a un dossier /disk ( html , css , js)



# Étape 2 : servir avec Nginx 
FROM nginx:alpine

# Copie le build React (Vite) généré à l'étape 1
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
