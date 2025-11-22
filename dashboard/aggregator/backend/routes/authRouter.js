import express from "express";
import jwt from "jsonwebtoken";
import { signupValidation, loginValidation } from "../middleware/AuthValidation.js";
import { signup, login } from "../controllers/AuthController.js";
import { UserModel } from "../models/user.js";

const router = express.Router();

// Signup
router.post("/signup", signupValidation, signup);

// Login
router.post("/login", loginValidation, login);

// Get logged-in user details
router.get("/me", async (req, res) => {
  try {
    const authHeader = req.headers.authorization;

    if (!authHeader)
      return res.json({ success: false, message: "No token provided" });

    const token = authHeader.split(" ")[1];
    if (!token)
      return res.json({ success: false, message: "Invalid token format" });

    const decoded = jwt.verify(token, process.env.JWT_SECRET);

    const user = await UserModel.findById(decoded.id).select("-password");
    if (!user)
      return res.json({ success: false, message: "User not found" });

    return res.json({ success: true, user });
  } catch (err) {
    console.log("ME Error:", err);
    return res.json({ success: false, message: "Invalid or expired token" });
  }
});

export default router;
