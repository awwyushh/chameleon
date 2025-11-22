// controllers/AttackController.js
import Attack from "../models/Attack.js";

// GET /attacks/dashboard
export const getDashboardData = async (req, res) => {
  try {
    const attacks = await Attack.find().sort({ createdAt: -1 }).limit(10);

    const total = await Attack.countDocuments();
    const sql = await Attack.countDocuments({ "classification.label": "SQLi" });
    const xss = await Attack.countDocuments({ "classification.label": "XSS" });
    const brute = await Attack.countDocuments({ "classification.label": "Bruteforce" });

    const avgDelayAgg = await Attack.aggregate([
      { $group: { _id: null, avgDelay: { $avg: "$tarpitDelayMs" } } }
    ]);

    const avgDelay = avgDelayAgg.length ? avgDelayAgg[0].avgDelay : 0;

    return res.json({
      success: true,
      kpis: {
        totalAttacks: total,
        sqlAttacks: sql,
        xssAttacks: xss,
        bruteAttacks: brute,
        avgDelay: Math.round(avgDelay),
      },
      recent: attacks,
    });
  } catch (err) {
    console.error("Dashboard error:", err);
    return res.status(500).json({ success: false, message: "Failed to load dashboard" });
  }
};
