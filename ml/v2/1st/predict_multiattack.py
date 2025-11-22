import pandas as pd
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.linear_model import LogisticRegression
from sklearn.svm import LinearSVC
from sklearn.model_selection import train_test_split
from sklearn.metrics import classification_report
import joblib

print("Training MULTI-ATTACK MODEL...")

# ============================
# 1. LOAD BOTH DATASETS
# ============================
df1 = pd.read_csv("data/SQLInjection_XSS_CommandInjection_MixDataset.1.0.0.csv")
df2 = pd.read_csv("data/SQLInjection_XSS_MixDataset.1.0.0.csv")

df = pd.concat([df1, df2], ignore_index=True)
print("Total rows after merge:", len(df))

# ============================
# 2. CLEANING
# ============================
df = df.drop_duplicates()
df = df.dropna(subset=["Sentence"])
df = df[df["Sentence"].str.strip() != ""]

# ============================
# 3. CREATE LABEL COLUMN
# ============================
def get_label(row):
    if row.get("SQLInjection", 0) == 1:
        return "SQLI"
    if row.get("XSS", 0) == 1:
        return "XSS"
    if row.get("CommandInjection", 0) == 1:
        return "CMDI"
    return "BENIGN"

df["Label"] = df.apply(get_label, axis=1)

# ============================
# 4. TAKE RANDOM 5000 SAMPLES
# ============================
df = df.sample(n=5000, random_state=42)
print("Using random 5000 rows.")

X = df["Sentence"]
y = df["Label"]

# ============================
# 5. VECTORIZATION
# ============================
vectorizer = TfidfVectorizer(
    analyzer="char",
    ngram_range=(2, 4),
    max_features=20000
)

X_vec = vectorizer.fit_transform(X)

# ============================
# 6. TRAIN / TEST SPLIT
# ============================
X_train, X_test, y_train, y_test = train_test_split(
    X_vec, y, test_size=0.2, random_state=42
)

# ============================
# 7. TRAIN MODELS
# ============================
print("Training Logistic Regression...")
lr = LogisticRegression(max_iter=3000)
lr.fit(X_train, y_train)

print("Training Linear SVM...")
svm = LinearSVC()
svm.fit(X_train, y_train)

# ============================
# 8. EVALUATION
# ============================
print("\n=== Logistic Regression Report ===")
print(classification_report(y_test, lr.predict(X_test)))

print("\n=== Linear SVC Report ===")
print(classification_report(y_test, svm.predict(X_test)))

# ============================
# 9. SAVE PKL MODELS
# ============================
joblib.dump(lr, "lr_multi.pkl")
joblib.dump(svm, "svm_multi.pkl")
joblib.dump(vectorizer, "vectorizer_multi.pkl")

print("\nPickle models saved successfully!")

# ============================
# 10. EXPORT LOGISTIC REGRESSION TO ONNX
# ============================
# ============================
# 10. EXPORT TO ONNX
# ============================
print("\nExporting ONNX model...")

from skl2onnx import convert_sklearn
from skl2onnx.common.data_types import StringTensorType, FloatTensorType
import onnx

# TF-IDF Vectorizer → ONNX
initial_type_text = [('input', StringTensorType([None, 1]))]

tfidf_onnx = convert_sklearn(
    vectorizer,
    'tfidf_vectorizer',
    initial_types=initial_type_text,
    target_opset={'': 21, 'ai.onnx.ml': 3}
)

with open("tfidf_vectorizer.onnx", "wb") as f:
    f.write(tfidf_onnx.SerializeToString())

print("✓ TF-IDF exported as tfidf_vectorizer.onnx")

# Logistic Regression → ONNX
initial_type = [('input', FloatTensorType([None, X_vec.shape[1]]))]

lr_onnx = convert_sklearn(
    lr,
    'logistic_regression',
    initial_types=initial_type,
    target_opset={'': 21, 'ai.onnx.ml': 3}
)

with open("logistic_multiattack.onnx", "wb") as f:
    f.write(lr_onnx.SerializeToString())

print("✓ Logistic Regression exported as logistic_multiattack.onnx")

# Linear SVM → ONNX (optional)
svm_onnx = convert_sklearn(
    svm,
    'linear_svm',
    initial_types=initial_type,
    target_opset={'': 21, 'ai.onnx.ml': 3}
)

with open("svm_multiattack.onnx", "wb") as f:
    f.write(svm_onnx.SerializeToString())

print("✓ SVM exported as svm_multiattack.onnx")
print("\n🎉 ONNX Export Completed Successfully!")
