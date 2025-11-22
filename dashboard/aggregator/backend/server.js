import express from "express";
import dotenv from "dotenv";
import { connectDB } from "./models/db.js";
import bodyParser from "body-parser";
import cors from "cors";
import attackRoutes from "./routes/AttackRoutes.js";




// Routes
import authRouter from "./routes/authRouter.js";


dotenv.config();
const app = express();
const PORT = process.env.PORT || 5000;

// Connect to MongoDB (same DB for users + projects)
connectDB();

// Middleware
app.use(bodyParser.json());
app.use(cors());

// Routes
app.use("/auth", authRouter);            // Login/Auth
app.use("/attacks", attackRoutes);

// Server start
app.listen(PORT, () => {
  console.log(`🚀 Server running on port ${PORT}`);
});
