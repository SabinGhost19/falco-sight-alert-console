# Etapa de build
FROM node:20-alpine AS builder

WORKDIR /app

# Instalăm dependențele
COPY package*.json ./
RUN npm ci

# Copiem codul și construim aplicația
COPY . .
RUN npm run build

# Etapa finală - un server web lightweight (Nginx)
FROM nginx:alpine

# Copiem config-ul pentru Vue/SPA rounting (Redirecționează către index.html)
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Adăugăm fișierele statice construite
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
