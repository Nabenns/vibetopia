#!/bin/bash
# VIBETOPIA Mac Fix — Register growtopia:// URL scheme untuk Growtopia 5.50+
# Run: chmod +x vibetopia-mac-fix.sh && ./vibetopia-mac-fix.sh

set -e

APP_PATH="/Applications/Growtopia.app"

if [ ! -d "$APP_PATH" ]; then
    echo "❌ Growtopia.app not found. Drag your Growtopia.app path:"
    read -r APP_PATH
    if [ ! -d "$APP_PATH" ]; then
        echo "❌ Still not found. Aborting."
        exit 1
    fi
fi

PLIST="$APP_PATH/Contents/Info.plist"

echo "🔧 Patching $PLIST ..."

# Backup
cp "$PLIST" "$PLIST.bak"
echo "   ✅ Backup: $PLIST.bak"

# Check if CFBundleURLTypes already exists
if /usr/libexec/PlistBuddy -c "Print :CFBundleURLTypes" "$PLIST" &>/dev/null; then
    echo "   ⚠️  CFBundleURLTypes already exists — checking for growtopia://"
    if /usr/libexec/PlistBuddy -c "Print :CFBundleURLTypes:0:CFBundleURLSchemes" "$PLIST" 2>/dev/null | grep -q "growtopia"; then
        echo "   ✅ growtopia:// already registered. Nothing to do."
        exit 0
    fi
else
    # Add CFBundleURLTypes array
    /usr/libexec/PlistBuddy -c "Add :CFBundleURLTypes array" "$PLIST"
    echo "   ✅ Added CFBundleURLTypes array"
fi

# Add growtopia:// URL scheme
/usr/libexec/PlistBuddy -c "Add :CFBundleURLTypes:0 dict" "$PLIST" 2>/dev/null || true
/usr/libexec/PlistBuddy -c "Add :CFBundleURLTypes:0:CFBundleURLName string 'com.ubisoft.growtopia.url'" "$PLIST" 2>/dev/null || true
/usr/libexec/PlistBuddy -c "Add :CFBundleURLTypes:0:CFBundleURLSchemes array" "$PLIST" 2>/dev/null || true
/usr/libexec/PlistBuddy -c "Add :CFBundleURLTypes:0:CFBundleURLSchemes:0 string 'growtopia'" "$PLIST" 2>/dev/null || true

echo "   ✅ growtopia:// URL scheme registered"

# Re-sign the app (ad-hoc signing — works for local dev)
echo "🔐 Re-signing $APP_PATH ..."
codesign --force --deep --sign - "$APP_PATH" 2>&1 || echo "   ⚠️  Re-sign failed (SIP may block). You may need to: sudo spctl --master-disable"

echo ""
echo "🎸 VIBETOPIA Mac fix applied!"
echo ""
echo "Next steps:"
echo "  1. sudo spctl --master-disable  (if Gatekeeper blocks)"
echo "  2. Open Growtopia.app"
echo "  3. Edit /private/etc/hosts (sudo nano /private/etc/hosts):"
echo "     103.253.213.178  gtps.bensserver.cloud"
echo "     103.253.213.178  cdn.bensserver.cloud"
echo "  4. Login with: Nabenns / qwerty77"
echo ""
echo "To revert: cp $PLIST.bak $PLIST"
