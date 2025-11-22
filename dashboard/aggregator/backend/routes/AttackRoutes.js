// backend/routes/AttackRoutes.js
import express from "express";
import Attack from "../models/Attack.js";
import { getDashboardData } from "../controllers/AttackController.js";

const router = express.Router();

router.get("/dashboard", getDashboardData);

router.post("/", async (req, res) => {
  try {
    const attack = new Attack(req.body);
    await attack.save();
    return res.json({ success: true, attack });
  } catch (err) {
    console.error("POST /attacks error:", err);
    return res.status(500).json({ success: false, message: "Failed to insert attack" });
  }
});

export default router;
