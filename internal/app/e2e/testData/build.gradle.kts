plugins {
    kotlin("jvm") version "2.1.10"
}

repositories {
    mavenCentral()
}

sourceSets {
    main {
        kotlin.setSrcDirs(listOf("cases"))
    }
}

dependencies {
    implementation(platform("io.ktor:ktor-bom:3.1.1"))
    implementation("io.ktor:ktor-client-core")
    implementation("io.ktor:ktor-client-cio")
    implementation("ch.qos.logback:logback-classic:1.4.14")
}