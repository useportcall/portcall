<#import "template.ftl" as layout>
<@layout.emailLayout>
<h1>Reset your password</h1>
<p>We received a request to reset the password for your <strong style="color: #0a0a0a;">${realmName}</strong> account.</p>
<p>Click the button below to choose a new password:</p>
<p style="text-align: center; padding: 8px 0;">
    <a href="${link}" class="button">Reset password</a>
</p>
<p>Or copy and paste this link into your browser:</p>
<p class="link-box">
    <a href="${link}">${link}</a>
</p>
<p>This link expires in <strong style="color: #0a0a0a;">${linkExpirationFormatter(linkExpiration)}</strong>.</p>
<p class="hint">If you didn't request this, you can safely ignore this email.</p>
</@layout.emailLayout>
